package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/primawk/gator-go/internal/database"
	"html"
	"io"
	"net/http"
	"time"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	resource := &RSSFeed{}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, feedURL, nil)

	if err != nil {
		return nil, fmt.Errorf("Error creating a fetch request: %s", err)
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{Timeout: 10 * time.Second}

	res, err := client.Do(req)

	body, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}

	if err != nil {
		return nil, err
	}

	if err := xml.Unmarshal(body, &resource); err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %v", err)
	}

	resource.Channel.Title = html.UnescapeString(resource.Channel.Title)
	resource.Channel.Description = html.UnescapeString(resource.Channel.Description)

	for i, rssItem := range resource.Channel.Item {
		resource.Channel.Item[i].Title = html.UnescapeString(rssItem.Title)
		resource.Channel.Item[i].Description = html.UnescapeString(rssItem.Description)
	}

	return resource, nil
}

func handlerFetchFeed(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("time between request is required.")
	}

	timeBetweenRequests, _ := time.ParseDuration(cmd.args[0])

	fmt.Printf("Collecting feeds every %s\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s, cmd)
	}

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	payload := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), payload)

	if err != nil {
		return fmt.Errorf("error saving a feed: %s\n", err)
	}

	payloadCreateFeed := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, errCreateFeed := s.db.CreateFeedFollow(context.Background(), payloadCreateFeed)

	if err != nil {
		return fmt.Errorf("error creating a feed follow: %s", errCreateFeed)
	}

	fmt.Println(feed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())

	if err != nil {
		return fmt.Errorf("error fetching feeds: %s\n", err)
	}

	fmt.Println(feeds)

	return nil
}

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	feedId, err := s.db.GetFeedId(context.Background(), cmd.args[0])

	if err != nil {
		return fmt.Errorf("error getting the feed id: %s", err)
	}

	payload := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedId,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), payload)

	if err != nil {
		return fmt.Errorf("error creating a feed follow: %s", err)
	}

	fmt.Println(feedFollow)

	return nil
}

func handlerFetchFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)

	if err != nil {
		return fmt.Errorf("error getting the following feeds: %s", err)
	}

	fmt.Println(feeds)

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return errors.New("feed url is required.")
	}

	feedId, err := s.db.GetFeedId(context.Background(), cmd.args[0])

	if err != nil {
		return fmt.Errorf("error getting the feed details: %s", err)
	}

	payload := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feedId,
	}

	if errDeleteFeedFollow := s.db.DeleteFeedFollow(context.Background(), payload); errDeleteFeedFollow != nil {
		return fmt.Errorf("error to unfollow a feed: %s", err)
	}

	return nil
}

func scrapeFeeds(s *state, cmd command) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())

	if err != nil {
		return fmt.Errorf("error getting the next feed to fetch: %s", err)
	}

	payload := database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true, // true means it's not NULL
		},
		ID: feed.ID,
	}

	if err := s.db.MarkFeedFetched(context.Background(), payload); err != nil {
		return fmt.Errorf("error marking the feed: %s", err)
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)

	if err != nil {
		return fmt.Errorf("Error fetching a feed: %s", err)
	}

	for _, item := range rssFeed.Channel.Item {
		fmt.Println(item.Title)
	}

	return nil
}
