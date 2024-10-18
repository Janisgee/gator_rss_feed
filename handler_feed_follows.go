package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Janisgee/gator_rss_feed/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("url is not provided from argument")
	}

	feedUrl := cmd.Args[0]
	username := s.cfg.CurrentUserName

	// // Create an empty context
	ctx := context.Background()

	// Get feed
	feed, err := s.db.GetFeedByUrl(ctx, feedUrl)
	if err != nil {
		return fmt.Errorf("error in getting feed from database with provided URL(%s):%w", feedUrl, err)
	}

	// Prepare the parameters
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
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

func handlerListFeedFollows(s *state, cmd command, user database.User) error {
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

func handleUnfollow(s *state, cmd command, user database.User) error {
	// Create an empty context
	ctx := context.Background()

	if len(cmd.Args) != 1 {
		return fmt.Errorf("error in missing argument. feed url is needed for unfollowing feed action")
	}

	feedUrl := cmd.Args[0]

	// Get feed
	feed, err := s.db.GetFeedByUrl(ctx, feedUrl)
	if err != nil {
		return fmt.Errorf("error in getting feed from database with provided URL(%s):%w", feedUrl, err)
	}

	//params for unfollow feed
	params := database.DeleteFeedFollowsForUserParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}
	// Unfollow feed by user (userID & feedUrl)
	deletedRows, err := s.db.DeleteFeedFollowsForUser(ctx, params)
	if err != nil {
		return fmt.Errorf("error in deleting feed(url:%s) follows for user(%s) and :%w", feedUrl, user.Name, err)
	}

	fmt.Printf("Unfollowed successfully:\n")
	fmt.Printf("* Feed Name: %s Feed URL: %s\n", feed.Name, feedUrl)
	fmt.Printf("* Feed Info: %s\n", deletedRows)

	return nil
}
