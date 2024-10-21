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
	//Login user by (input)username
	appCommands.Register("login", handlerLogin)
	// Register new user by (input)username
	appCommands.Register("register", handlerRegister)
	// Delete all user database
	appCommands.Register("reset", handlerReset)
	// Get all users from database
	appCommands.Register("users", handlerUsers)
	// Fetch feed By (input) time between requests
	appCommands.Register("agg", handlerAgg)
	// All feed from logined user by (input) feedName and feedUrl
	appCommands.Register("addfeed", middlewareLoggedIn(handlerAddFeed))
	//List all the feeds from database
	appCommands.Register("feeds", handlerListFeeds)
	//Get feed_name, username by (input) FeedURL
	appCommands.Register("follow", middlewareLoggedIn(handlerFollow))
	//Get feed_name by (input) LOGIN current user
	appCommands.Register("following", middlewareLoggedIn(handlerListFeedFollows))
	//Delete feedfollow by (input) userID and feedUrl
	appCommands.Register("unfollow", middlewareLoggedIn(handleUnfollow))
	// Browse posts by (input) limited number of post/ default 2 posts
	appCommands.Register("browse", middlewareLoggedIn(handlerBrowse))

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

// Terminal command for development:

//psql "postgres://postgres:postgres@localhost:5432/gator"
//goose postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up

// using psql to find your newly created users table
//psql gator
//\dt

//sqlc generate
//go build -o gator_rss_feed
