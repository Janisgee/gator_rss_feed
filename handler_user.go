package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Janisgee/gator_rss_feed/internal/database"
	"github.com/google/uuid"
)

// Signature of all command handlers
func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("a username is required in arguments")
	}

	// Check if the given username exist in database

	// Get username from argument
	username := cmd.Args[0]
	// Create an empty context
	ctx := context.Background()

	// Check if queries username is in database
	_, err := s.db.GetUser(ctx, username)
	if err != nil {
		fmt.Println("Given username does not exist in database")
		os.Exit(1)
	}

	// Set the user to the given username
	err = s.cfg.SetUser(username)
	if err != nil {
		return errors.New("error with setting current user")
	}
	fmt.Printf("User switched successfully: %s\n", cmd.Args[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("a name is required in arguments")
	}
	// Get username from argument
	username := cmd.Args[0]

	// Create an empty context
	ctx := context.Background()

	// Check if queries username is in database
	_, err := s.db.GetUser(ctx, username)
	if err == nil {
		fmt.Println("already have the same username saved in database as the input username")
		os.Exit(1)
	}

	// Create user if username is not exist in database
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	}
	userData, err := s.db.CreateUser(ctx, params)
	if err != nil {
		return errors.New("error in creating user into database")
	}

	err = s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("The user created successfully.\n New user data:%v\n", userData)

	return nil
}

func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	// Delete all users in database
	err := s.db.DeleteAllUsers(ctx)
	if err != nil {
		return fmt.Errorf("error deleting all users:%w", err)
	}

	fmt.Print("All users have been deleted successfully!")

	return nil

}

func handlerUsers(s *state, cmd command) error {
	ctx := context.Background()
	// Delete all users in database
	allUsers, err := s.db.GetUsers(ctx)
	if err != nil {
		return fmt.Errorf("error in getting all users:%w", err)
	}

	for i := 0; i < len(allUsers); i++ {
		if allUsers[i] == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", allUsers[i])
		} else {
			fmt.Printf("* %s \n", allUsers[i])
		}
	}
	return nil
}

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
