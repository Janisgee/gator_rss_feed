package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Janisgee/gator_rss_feed/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("url is not provided from argument")
	}

	feedUrl := cmd.Args[0]
	username := s.cfg.CurrentUserName

	// Create an empty context
	ctx := context.Background()

	// Get userID
	userInfo, err := s.db.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("current user is not regist yet:%w", err)
	}

	// Get feed
	feed, err := s.db.GetFeed(ctx, feedUrl)
	if err != nil {
		return fmt.Errorf("error in getting feed from database with provided URL(%s):%w", feedUrl, err)
	}

	// Prepare the parameters
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userInfo.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(ctx, params)
	if err != nil {
		return fmt.Errorf("error in creating feed with feedURL(%s) from user(%s):%w", feedUrl, username, err)
	}

	fmt.Println("Successfully created feed:")
	fmt.Printf("Feed Name:         %s\n", feed.Name)
	fmt.Printf("Current username:  %s\n", username)
	fmt.Println()

	return nil
}

func handlerListFeedFollows(s *state, cmd command) error {
	currentUser := s.cfg.CurrentUserName

	// Create an empty context
	ctx := context.Background()

	// Get feeds follow By user
	feedFollow, err := s.db.GetFeedFollowsForUser(ctx, currentUser)
	if err != nil {
		return fmt.Errorf("error in getting feed follow from current user(%s):%w", currentUser, err)
	}

	if len(feedFollow) == 0 {
		fmt.Printf("cannot found any feed follow by current user(%s)\n", currentUser)
		return nil
	}
	// Print all the names of the feeds the current user is following
	fmt.Printf("Names of the feeds %s following:\n", currentUser)
	for i := range feedFollow {
		fmt.Printf("* %s\n", feedFollow[i].FeedName)
	}
	fmt.Println()
	return nil
}
