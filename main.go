package main

import (
	"fmt"
	"github.com/nunseik/gator/internal/config"
	"os"
	"log"
)


func main() {
	config, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
	}
	newState := &state{config: &config}
	commands := &commands{commands: make(map[string]func(*state, command) error)}
	commands.register("login", handlerLogin)
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

	err = config.SetUser()
	if err != nil {
		log.Fatalf("couldn't set current user: %v", err)
	}

}