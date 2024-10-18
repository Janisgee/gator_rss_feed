package main

import (
	"context"
	"fmt"

	"github.com/Janisgee/gator_rss_feed/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, cmd command) error {

		username := s.cfg.CurrentUserName

		// Create an empty context
		ctx := context.Background()

		// Get userID
		user, err := s.db.GetUser(ctx, username)
		if err != nil {
			return fmt.Errorf("current user is not regist yet:%w", err)
		}
		return handler(s, cmd, user)
	}

}
