package main

import (
	"context"
	"fmt"
	"time"

	"github.com/atmetz/rssgator/internal/database"
	"github.com/google/uuid"
)

// addfeed command
func handlerAddFeed(s *state, cmd command, user database.User) error {

	// Verify the correct number of arguments
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %s <name> <RSS feed url>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	// Get RSS Feed and add to DB
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %v", err)
	}

	// Automatically follow newly created feeds by current user
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed, user)
	fmt.Println()
	fmt.Println("Feed followed successfully:")
	printFollowing(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("=====================================")

	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
	fmt.Printf("* User:          %s\n", user.Name)
}

// feeds command
func handlerFeeds(s *state, cmd command) error {

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not get feeds: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	// Print all feeds

	fmt.Printf("Found %d feeds:\n", len(feeds))
	for _, feed := range feeds {
		user, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("couldn't get user: %v", err)
		}
		printFeed(feed, user)
		fmt.Println("=====================================")
	}

	return nil
}

// follow command
func handlerFollow(s *state, cmd command, user database.User) error {

	// Verify the correct number of arguments
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <RSS feed url>", cmd.Name)
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("feed does not exist: %v", err)
	}

	// Create followed feed

	row, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return fmt.Errorf("could not follow feed at %v: %v", cmd.Args[0], err)
	}

	fmt.Println("Feed followed successfully:")
	printFollowing(row.UserName, row.FeedName)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

// following command
func handlerFollowing(s *state, cmd command, user database.User) error {

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error displaying following: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}

	// Print a list of feeds followed by current user

	fmt.Printf("%s is following:\n", s.cfg.CurrentUserName)
	for _, feed := range feeds {
		fmt.Printf("* %s\n", feed.FeedName)
	}
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

func printFollowing(username, feedname string) {

	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)

}

// unfollow command
func handlerUnfollow(s *state, cmd command, user database.User) error {

	// Verify the correct number of arguments
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <RSS feed url>", cmd.Name)
	}

	// find feed
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("feed does not exist: %v", err)
	}

	// unfollow based on userid and feedid
	err = s.db.UnfollowFeed(context.Background(), database.UnfollowFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return fmt.Errorf("could not unfollow: %v", err)
	}

	fmt.Println("Successfully unfollowed feed.")
	return nil
}
