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
