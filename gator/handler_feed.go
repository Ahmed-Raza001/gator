package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ahmed/gator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("usage: addfeed <name> <url>")
	}

	name := cmd.args[0]
	feedURL := cmd.args[1]

	now := time.Now()

	feed, err := s.db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			Name:      name,
			Url:       feedURL,
			UserID:    user.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("could not create feed: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("could not follow feed: %w", err)
	}

	fmt.Println("Feed created successfully!")
	fmt.Printf("Feed: %+v\n", feed)
	fmt.Printf(
		"%s is now following %s\n",
		feedFollow.UserName,
		feedFollow.FeedName,
	)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("could not get feeds: %w", err)
	}

	for _, feed := range feeds {
		fmt.Printf("* %s\n", feed.Name)
		fmt.Printf("  URL: %s\n", feed.Url)
		fmt.Printf("  User: %s\n", feed.UserName)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: follow <url>")
	}

	feedURL := cmd.args[0]

	feed, err := s.db.GetFeedByURL(
		context.Background(),
		feedURL,
	)
	if err != nil {
		return fmt.Errorf("could not get feed: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("could not create feed follow: %w", err)
	}

	fmt.Printf(
		"%s is now following %s\n",
		feedFollow.UserName,
		feedFollow.FeedName,
	)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(
		context.Background(),
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("could not get feed follows: %w", err)
	}

	fmt.Printf("Feeds followed by %s:\n", user.Name)

	for _, feedFollow := range feedFollows {
		fmt.Printf("* %s\n", feedFollow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: unfollow <url>")
	}

	feedURL := cmd.args[0]

	// Get feed by URL.
	feed, err := s.db.GetFeedByURL(
		context.Background(),
		feedURL,
	)
	if err != nil {
		return fmt.Errorf("could not get feed: %w", err)
	}

	// Delete feed follow record.
	err = s.db.DeleteFeedFollow(
		context.Background(),
		database.DeleteFeedFollowParams{
			UserID: user.ID,
			FeedID: feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("could not unfollow feed: %w", err)
	}

	fmt.Printf("%s has unfollowed %s\n", user.Name, feed.Name)

	return nil
}
