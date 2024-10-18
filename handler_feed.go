package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Janisgee/gator_rss_feed/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return errors.New("not enough arguments. arguments should contain name and url")
	}
	// Get arguments input
	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	// Create an empty context
	ctx := context.Background()

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

	fmt.Printf("The feed(%s) was created successfully:\n", feedName)
	printFeed(feedData)
	fmt.Println()
	fmt.Println("=================================================")

	if len(cmd.Args) >= 2 {
		cmd.Args = cmd.Args[1:]
	}

	err = handlerFollow(s, cmd, user)
	if err != nil {
		return err
	}

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:          %s\n", feed.ID)
	fmt.Printf("* CreatedAt:   %v\n", feed.CreatedAt)
	fmt.Printf("* UpdatedAt:   %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:        %s\n", feed.Name)
	fmt.Printf("* Url:         %s\n", feed.Url)
	fmt.Printf("* UserID:      %s\n", feed.UserID)
}

func handlerListFeeds(s *state, cmd command) error {
	// Create an empty context
	ctx := context.Background()

	allFeeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return fmt.Errorf("error in getting all feeds:%w", err)
	}

	if len(allFeeds) == 0 {
		fmt.Printf("No feeds found.%d feed in database.\n", len(allFeeds))
		return nil
	}

	fmt.Printf("Found %d feeds:\n", len(allFeeds))
	fmt.Println("=============================")
	fmt.Println()
	for i := 0; i < len(allFeeds); i++ {
		fmt.Printf("* Feed Name:  %s\n", allFeeds[i].FeedName)
		fmt.Printf("* Url:        %s\n", allFeeds[i].Url)
		fmt.Printf("* Username:   %s\n", allFeeds[i].UserName)
		fmt.Println()
		fmt.Println("=============================")
		fmt.Println()
	}

	return nil
}
