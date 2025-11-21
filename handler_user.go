package main

import (
	"context"
	"fmt"
	"time"

	"github.com/atmetz/rssgator/internal/database"
	"github.com/google/uuid"
)

// Login function
func handlerLogins(s *state, cmd command) error {
	//set context for queries
	con := context.Background()

	// Verify the correct number of arguments
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	// Check if user exists
	_, err := s.db.GetUser(con, cmd.Args[0])

	if err != nil {
		return fmt.Errorf("user does not exist: %v", err)
	}

	// Set current user to argument passed in
	err = s.cfg.SetUser(cmd.Args[0])

	if err != nil {
		return fmt.Errorf("couldn't set current user: %v", err)
	}

	fmt.Println("User switched successfully!")

	return nil
}

// Register user function
func handlerRegister(s *state, cmd command) error {
	//set context for queries
	con := context.Background()

	// Verify the correct number of arguments
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	// Check if user exists
	_, err := s.db.GetUser(con, cmd.Args[0])

	if err == nil {
		return fmt.Errorf("user already exists: %v", err)
	}

	// Create New User
	createdUser, err := s.db.CreateUser(con, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.Args[0],
	})

	if err != nil {
		return fmt.Errorf("error creating new user %s, %v", cmd.Args[0], err)
	}

	// Set current user to newly created user
	err = s.cfg.SetUser(createdUser.Name)

	if err != nil {
		return fmt.Errorf("couldn't set current user: %v", err)
	}

	fmt.Println("User created successfully:")
	printUser(createdUser)
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
