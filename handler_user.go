package main

import (
	"errors"
	"fmt"
)

// Signature of all command handlers
func HandlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("a username is required")
	}

	// Set the user to the given username
	if s.cfg != nil {
		s.cfg.SetUser(cmd.Args[0])
		fmt.Printf("The user has been set:%s\n", cmd.Args[0])
	} else {
		return errors.New("error with setting current user")
	}

	return nil
}
