package gatorapi

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
	"html")

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

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	client := &http.Client{
		Timeout: 10 * time.Second, // Set a timeout for the request
	}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to fetch feed: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return &RSSFeed{}, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to read response body: %w", err)
	}
	var feed RSSFeed
	err = xml.Unmarshal(resBody, &feed)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("failed to parse feed: %w", err)
	}

	for _, item := range feed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	return &feed, nil
}