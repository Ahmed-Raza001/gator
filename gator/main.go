package main

import (
	"database/sql"
	"fmt"
	"github.com/ahmed/gator/internal/database"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/ahmed/gator/internal/config"
)

// state holds the application state.
type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {

	// Read config file
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	// Create application state
	s := &state{
		db:  dbQueries,
		cfg: &cfg,
	}

	// Initialize commands
	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}

	// Register login command
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
    cmds.register("follow", middlewareLoggedIn(handlerFollow))
    cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("feeds", handlerFeeds)
	cmds.register("agg", handlerAgg)
	cmds.register("users", handlerUsers)
	cmds.register("reset", handlerReset)
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	// Check command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Error: Not enough arguments provided")
		os.Exit(1)
	}

	// Create command
	cmd := command{
		name: os.Args[1],
		args: os.Args[2:],
	}

	// Run command
	err = cmds.run(s, cmd)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

}
