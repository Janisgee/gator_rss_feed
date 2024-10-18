package main

import (
	"context"
	"fmt"
	"time"
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
		scrapeFeeds(s)
		fmt.Println()
		fmt.Println("==================================")
		fmt.Println()
	}

}

func scrapeFeeds(s *state) {
	ctx := context.Background()
	// Get next feed to fetch from DB
	nextFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		fmt.Printf("error in getting next feed to fetch\n")
	}

	// Mark it as fetched
	err = s.db.MarkFeedFetched(ctx, nextFeed.ID)
	if err != nil {
		fmt.Printf("error in marking next feed(feed name: %s) as fetch.\n", nextFeed.Name)
	}

	// Fetch the feed using URL
	rssFeed, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		fmt.Printf("error in fetching next feed(url:%s)\n", nextFeed.Url)
	}

	for _, items := range rssFeed.Channel.Item {
		fmt.Printf("RssFeed Item Title:  %s\n", items.Title)
	}
	fmt.Printf("Feed %s collected, %v posts found\n", nextFeed.Name, len(rssFeed.Channel.Item))

}
