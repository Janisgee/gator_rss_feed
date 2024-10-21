package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Janisgee/gator_rss_feed/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {

	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("error in missing or wrong argument input <time between requests eg. 1s, 1m, 1h, etc.>")
	}
	// Parse input argument into a time.Duration value
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error in parseing duration from input:%s :%w", cmd.Args[0], err)
	}
	if timeBetweenRequests < 0 {
		return fmt.Errorf("time between requests must be greater than 0")
	}

	fmt.Printf("Collecting feeds every %s...\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
		fmt.Println()
		fmt.Println("==================================")
		fmt.Println()
	}

}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	// Get next feed to fetch from DB
	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		fmt.Printf("error in getting next feed to fetch\n")
	}

	// Mark it as fetched
	_, err = s.db.MarkFeedFetched(ctx, nextFeed.ID)
	if err != nil {
		fmt.Printf("error in marking next feed(feed name: %s) as fetch.\n", nextFeed.Name)
	}

	// Fetch the feed using URL
	rssFeed, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		fmt.Printf("error in fetching next feed(url:%s)\n", nextFeed.Url)
	}

	for _, items := range rssFeed.Channel.Item {
		// Parse time into time
		publishedAt := sql.NullTime{}
		t, err := time.Parse(time.RFC1123Z, items.PubDate)
		if err != nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		//Create posts table
		params := database.CreatePostsParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     items.Title,
			Url:       items.Link,
			Description: sql.NullString{
				String: items.Description,
				Valid:  true,
			},
			PublishedAt: publishedAt,
			FeedID:      nextFeed.ID,
		}
		_, err = s.db.CreatePosts(ctx, params)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			fmt.Printf("Couldn't create post: %v\n", err)
			continue
		}

	}
	fmt.Printf("Feed %s collected, %v posts found", nextFeed.Name, len(rssFeed.Channel.Item))

	return nil
}
