package main

import (
	"errors"
	"fmt"

	"github.com/nunseik/gator/internal/config"
	"github.com/nunseik/gator/internal/database"
)

type state struct {
	db *database.Queries
	config *config.Config
}

type command struct {
	name     string
	commands []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.commands) == 0 {
		return errors.New("login command requires a username")
	}
	s.config.CurrentUserName = cmd.commands[0]
	fmt.Println("User set successfully:", s.config.CurrentUserName)
	return nil
}

type commands struct {
	commands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if handler, exists := c.commands[cmd.name]; exists {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.name)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commands[name] = f
}
