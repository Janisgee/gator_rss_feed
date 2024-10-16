package main

import (
	"fmt"
	"os"

	"github.com/Janisgee/gator_rss_feed/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config::%v\n", err)
	}

	fmt.Println("Read old Config:", cfg)

	appState := state{
		cfg: &cfg,
	}

	appCommands := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	appCommands.Register("login", HandlerLogin)

	args := os.Args
	if len(args) < 2 {
		fmt.Println("Error: not enough arguments were provided.")
		os.Exit(1)
	}
	// Split the command-line arguments into command name and the arguments slice
	commandName := args[1]
	commandArgs := args[2:]

	appCommand := command{
		Name: commandName,
		Args: commandArgs,
	}

	err = appCommands.Run(&appState, appCommand)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

}
