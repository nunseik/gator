package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nunseik/gator/internal/config"
	"github.com/nunseik/gator/internal/database"
)

type state struct {
	db     *database.Queries
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
	// Check if the user exists
	_, err := s.db.GetUser(context.Background(), cmd.commands[0])
	if err != nil {
		return fmt.Errorf("could not find user: %v", err)
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

func handlerRegister(s *state, cmd command) error {
	if len(cmd.commands) == 0 {
		return errors.New("register command requires a username")
	}
	username := cmd.commands[0]
	// Generate a unique ID for the user using UUID into int32 somehow, not working with int32 directly
	userUniqueID := uuid.New()
	// Check if the user already exists
	existingUser, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		return fmt.Errorf("user %s already exists with ID %v", existingUser.Name, existingUser.ID)
	}
	newUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: userUniqueID, CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: username})
	if err != nil {
		return fmt.Errorf("could not create user: %v", err)
	}
	s.config.CurrentUserName = newUser.Name
	fmt.Println("User registered successfully:", newUser.Name, "/n with ID:", newUser.ID)
	return nil
}
func resetUsers(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {		
		return fmt.Errorf("could not reset users: %v", err)
	}
	return nil
}