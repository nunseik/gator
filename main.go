package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/nunseik/gator/internal/config"
	"github.com/nunseik/gator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config: ", err)
	}
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		fmt.Println("couldn't connect to database: ", err)
	}
	defer db.Close()
	dbQueries := database.New(db)
	newState := &state{config: &cfg}
	newState.db = dbQueries
	commands := &commands{commands: make(map[string]func(*state, command) error)}
	// Register commands
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", resetUsers)
	commands.register("users", getAllUsers)
	commands.register("agg", fetchCommand)
	commands.register("addfeed", createFeed)
	commands.register("feeds", getFeed)
	commands.register("follow", createFeedFollow)
	commands.register("following", getFeedFollows)
	// Add more commands as needed
	args := os.Args
	if len(args) < 2 {
		fmt.Println("No command provided")
		os.Exit(1)
	}
	cmd := command{name: args[1], commands: args[2:]}
	if err := commands.run(newState, cmd); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	err = cfg.SetUser()
	if err != nil {
		log.Fatalf("couldn't set current user: %v", err)
	}

}
