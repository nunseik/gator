package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nunseik/gator/internal/config"
	"github.com/nunseik/gator/internal/database"
	"github.com/nunseik/gator/internal/gatorapi"
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

func getAllUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not retrieve users: %v", err)
	}
	for user := range users {
		if users[user].Name == s.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", users[user].Name)
		} else {
			fmt.Printf("* %s\n", users[user].Name)
		}
	}
	return nil
}

func fetchCommand(s *state, cmd command) error {
	// if len(cmd.commands) == 0 {
	// 	return errors.New("fetch command requires a feed URL")
	// }
	// feedURL := cmd.commands[0]
	feedURL := "https://www.wagslane.dev/index.xml" // Example feed URL
	ctx := context.Background()
	feed, err := gatorapi.FetchFeed(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("could not fetch feed: %v", err)
	}
	fmt.Print(feed)
	return nil
}

func createFeed(s *state, cmd command, user database.User) error {
	if len(cmd.commands) < 2 {
		return errors.New("create command requires a feed name and URL")
	}
	feedName := cmd.commands[0]
	feedURL := cmd.commands[1]
	feedID := uuid.New()
	newFeed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        feedID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed: %v", err)
	}
	fmt.Print(newFeed)

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    newFeed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed follow: %v", err)
	}
	fmt.Printf("Feed created successfully: %s with URL: %s\n", newFeed.Name, newFeed.Url)
	return nil
}

func getFeed(s *state, cmd command) error {
	feed, err := s.db.GetFeed(context.Background())
	if err != nil {
		return fmt.Errorf("could not retrieve feed: %v", err)
	}
	for _, f := range feed {
		user, err := s.db.GetUserById(context.Background(), f.UserID)
		if err != nil {
			return fmt.Errorf("could not find user for feed %s: %v", f.Name, err)
		}
		fmt.Printf("Feed Name: %s, URL: %s, User Name: %s\n", f.Name, f.Url, user.Name)
	}
	if len(feed) == 0 {
		fmt.Println("No feeds found.")
	}

	return nil
}

func createFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.commands) < 1 {
		return errors.New("follow requires a URL")
	}
	feedURL := cmd.commands[0]
	feed, err := s.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("could not find feed with URL %s: %v", feedURL, err)
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("could not create feed follow: %v", err)
	}
	fmt.Printf("Successfully followed feed %s\n", feed.Name)
	return nil
}

func getFeedFollows(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("could not retrieve feed follows: %v", err)
	}
	for _, follow := range follows {
		feed, err := s.db.GetFeedById(context.Background(), follow.FeedID)
		if err != nil {
			return fmt.Errorf("could not find feed for follow %s: %v", follow.ID, err)
		}
		fmt.Printf("Follow ID: %s, Feed Name: %s\n", follow.ID, feed.Name)
	}
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("user not logged in: %v", err)
		}
		return handler(s, cmd, user)
	}
}