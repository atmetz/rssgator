package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/atmetz/rssgator/internal/database"
)

func handlerAgg(s *state, cmd command) error {

	// Verify the correct number of arguments
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name> <time duration, etc 1m>", cmd.Name)
	}

	time_between_requests := cmd.Args[0]

	timeBetweenRequests, err := time.ParseDuration(time_between_requests)

	if err != nil {
		return fmt.Errorf("invalide duration: %v", err)
	}

	fmt.Printf("Collecting feeds every %v...\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

func scrapeFeeds(s *state) {

	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())

	if err != nil {
		log.Println("cannot load next feed", err)
	}
	log.Println("Found a feed to fetch")
	scrapeFeed(s.db, nextFeed)

}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("couldn't mark feed %s fecthed: %v", feed.Name, err)
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("can not get feed %s: %v", feed.Url, err)
	}

	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}
