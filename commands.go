package main

import (
	"errors"
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

// This method registers a new handler function for a command name
func (c *commands) Register(name string, f func(*state, command) error) {
	if name != "" && f != nil {
		c.registeredCommands[name] = f
	}
}

// This method runs a given command with the provided state if it exists
func (c *commands) Run(s *state, cmd command) error {
	if s.cfg == nil {
		return errors.New("error in no provided state when running a given command")
	}
	handlerFunc, exists := c.registeredCommands[cmd.Name]
	if !exists {
		fmt.Println("No handler found:", cmd.Name)
		return errors.New("error in no handler found when runing command")
	}
	return handlerFunc(s, cmd)

}
