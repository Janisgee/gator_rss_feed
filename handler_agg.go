package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Janisgee/gator_rss_feed/internal/database"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()
	feedURL := "https://www.wagslane.dev/index.xml"

	rssf, err := fetchFeed(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("error in fetching the feed at %s: %w", feedURL, err)
	}
	fmt.Printf("Result struct: %v\n", rssf)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		return errors.New("not enough arguments. arguments should contain name and url")
	}
	// Get arguments input
	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	// Create an empty context
	ctx := context.Background()
	// Get user ID
	username := s.cfg.CurrentUserName
	user, err := s.db.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("cannot find username from user database:%w", err)
	}

	//Create feeds table
	params := database.CreateFeedsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}
	feedData, err := s.db.CreateFeeds(ctx, params)
	if err != nil {
		return fmt.Errorf("error in creating feed(%s) in database:%w", feedName, err)
	}

	fmt.Printf("The feed(%s) was created successfully.\n New feed data:%v\n", feedName, feedData)

	return nil
}
