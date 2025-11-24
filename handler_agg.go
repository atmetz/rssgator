package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/atmetz/rssgator/internal/database"
	"github.com/google/uuid"
)

// agg command
func handlerAgg(s *state, cmd command) error {

	// Verify the correct number of arguments
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name> <time duration, etc 1m>", cmd.Name)
	}

	// set time interval between feed scrape
	time_between_requests := cmd.Args[0]

	timeBetweenRequests, err := time.ParseDuration(time_between_requests)

	if err != nil {
		return fmt.Errorf("invalid duration: %v", err)
	}

	fmt.Printf("Collecting feeds every %v...\n", timeBetweenRequests)

	//
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

// get new posts feeds
func scrapeFeeds(s *state) {

	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())

	if err != nil {
		log.Println("cannot load next feed", err)
	}
	log.Println("Found a feed to fetch")
	scrapeFeed(s.db, nextFeed)

}

// get new posts feeds
func scrapeFeed(db *database.Queries, feed database.Feed) {
	// mark feed as fetched
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("couldn't mark feed %s fecthed: %v", feed.Name, err)
	}

	// fetch feed
	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("can not get feed %s: %v", feed.Url, err)
	}

	// set different time format layouts
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
	}

	// set publishedAt with time format based on above list
	for _, item := range rssFeed.Channel.Item {

		publishedAt := sql.NullTime{}
		for _, l := range layouts {
			if t, err := time.Parse(l, item.PubDate); err == nil {
				publishedAt = sql.NullTime{
					Time:  t,
					Valid: true,
				}
			}
		}

		//fmt.Printf("Found post: %s\n", item.Title)
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				// duplicate key — ignore
				continue
			} else {
				log.Printf("create post error: %v", err)
				continue
			}
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))

}

// browse command
func handlerBrowse(s *state, cmd command, user database.User) error {

	if len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %s <limit, default 2>", cmd.Name)
	}

	// get post limit
	var limit int32
	limit = 2

	if len(cmd.Args) == 1 {
		_, err := fmt.Sscanf(cmd.Args[0], "%d", &limit)
		if err != nil {
			limit = 2
		}
	}

	// get number of posts based on limit
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})

	if err != nil {
		return fmt.Errorf("cannot browse posts: %v", err)
	}

	// print posts
	fmt.Printf("%s's most recent followed posts:\n", user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description.String)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}
