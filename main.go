package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Janisgee/gator_rss_feed/internal/config"
	"github.com/Janisgee/gator_rss_feed/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config::%v\n", err)
	}

	fmt.Println("Read old Config:", cfg)

	// Load database URL to config struct and sql.Open() a connection to my database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer db.Close()
	// Create a new *database.Queries, store it in state struct
	dbQueries := database.New(db)

	// Create a state instance
	appState := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	appCommands := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	appCommands.Register("login", handlerLogin)
	appCommands.Register("register", handlerRegister)
	appCommands.Register("reset", handlerReset)
	appCommands.Register("users", handlerUsers)
	appCommands.Register("agg", handlerAgg)

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

	err = appCommands.Run(appState, appCommand)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

}
