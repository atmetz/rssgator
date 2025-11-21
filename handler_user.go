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

	// Verify the correct number of arguments
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	// Check if user exists
	_, err := s.db.GetUser(context.Background(), cmd.Args[0])

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

	// Verify the correct number of arguments
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	// Check if user exists
	_, err := s.db.GetUser(context.Background(), cmd.Args[0])

	if err == nil {
		return fmt.Errorf("user already exists: %v", err)
	}

	// Create New User
	createdUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
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

func handlerReset(s *state, cmd command) error {

	//Reset users table

	err := s.db.Reset(context.Background())

	if err != nil {
		return fmt.Errorf("could not delete users: %v", err)
	}

	fmt.Println("database reset successfully!")
	return nil

}

func handlerUsers(s *state, cmd command) error {

	// print list of users and tag current user

	users, err := s.db.GetUsers(context.Background())

	if err != nil {
		return fmt.Errorf("could not get list of users: %v", err)
	}

	for _, user := range users {
		fmt.Printf(" * %s", user.Name)
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf((" (current)"))
		}
		fmt.Println()
	}

	return nil
}
