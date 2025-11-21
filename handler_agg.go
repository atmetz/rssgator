package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	feed := "https://www.wagslane.dev/index.xml"
	rssFeed, err := fetchFeed(context.Background(), feed)
	if err != nil {
		return fmt.Errorf("can not get feed from %s: %v", feed, err)
	}

	fmt.Println(rssFeed)

	return nil
}
