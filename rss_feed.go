package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Set("User-Agent", "gator")

	res, err := httpClient.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	feed := RSSFeed{}
	err = xml.Unmarshal(dat, &feed)
	if err != nil {
		return &RSSFeed{}, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, value := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(value.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(value.Description)
	}

	return &feed, nil
}
