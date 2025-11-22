package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	feed := "https://www.wagslane.dev/index.xml"
	// RSSfeed of feed variable
	rssFeed, err := fetchFeed(context.Background(), feed)
	if err != nil {
		return fmt.Errorf("can not get feed from %s: %v", feed, err)
	}

	fmt.Printf("Feed: %+v\n", rssFeed)

	return nil

}
