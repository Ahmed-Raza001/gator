package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
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
	// Create HTTP request.
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		feedURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	// Set User-Agent header.
	req.Header.Set("User-Agent", "gator")

	// Create HTTP client.
	client := &http.Client{}

	// Send request.
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not fetch feed: %w", err)
	}
	defer res.Body.Close()

	// Check response status.
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	// Read response body.
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	// Convert XML into RSSFeed struct.
	feed := &RSSFeed{}

	err = xml.Unmarshal(data, feed)
	if err != nil {
		return nil, fmt.Errorf("could not parse XML: %w", err)
	}

	// Decode HTML entities in channel fields.
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	// Decode HTML entities in item fields.
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(
			feed.Channel.Item[i].Title,
		)

		feed.Channel.Item[i].Description = html.UnescapeString(
			feed.Channel.Item[i].Description,
		)
	}

	return feed, nil
}
