package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ahmed/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			fmt.Printf("Error scraping feeds: %v\n", err)
		}
	}
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	// Get next feed.
	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return fmt.Errorf("could not get next feed: %w", err)
	}

	// Mark feed as fetched.
	err = s.db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		return fmt.Errorf("could not mark feed as fetched: %w", err)
	}

	fmt.Printf("Fetching feed: %s\n", feed.Name)

	// Fetch RSS feed.
	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("could not fetch feed: %w", err)
	}

	// Save posts.
	for _, item := range rssFeed.Channel.Item {

		publishedAt, err := parsePublishedAt(item.PubDate)

		if err != nil {
			publishedAt = sql.NullTime{}
		}

		_, err = s.db.CreatePost(ctx, database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})

		if err != nil {
			// Ignore duplicate URLs.
			var pqErr *pq.Error

			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				continue
			}

			fmt.Printf("Error saving post: %v\n", err)
			continue
		}

		fmt.Printf("Saved post: %s\n", item.Title)
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {

	// Default limit.
	limit := int32(2)

	// Optional limit argument.
	if len(cmd.args) > 1 {
		return fmt.Errorf("usage: browse [limit]")
	}

	if len(cmd.args) == 1 {
		parsedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil || parsedLimit <= 0 {
			return fmt.Errorf("limit must be a positive number")
		}

		limit = int32(parsedLimit)
	}

	// Get posts for current user.
	posts, err := s.db.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{
			UserID: user.ID,
			Limit:  limit,
		},
	)

	if err != nil {
		return fmt.Errorf("could not get posts: %w", err)
	}

	for _, post := range posts {
		fmt.Printf("\nTitle: %s\n", post.Title)
		fmt.Printf("URL: %s\n", post.Url)
		fmt.Printf("Description: %s\n", post.Description.String)
		fmt.Println("--------------------------------")
	}

	return nil
}

func parsePublishedAt(date string) (sql.NullTime, error) {
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"2006-01-02T15:04:05Z",
	}

	for _, layout := range layouts {
		parsedTime, err := time.Parse(layout, date)
		if err == nil {
			return sql.NullTime{
				Time:  parsedTime,
				Valid: true,
			}, nil
		}
	}

	return sql.NullTime{}, fmt.Errorf(
		"could not parse published date: %s",
		date,
	)
}
