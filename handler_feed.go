package main

import (
	"context"
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
	// if len(cmd.args) == 0 {
	// 	return errors.New("url is required.")
	// }

	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")

	if err != nil {
		return fmt.Errorf("Error fetching a feed: %s", err)
	}

	fmt.Println(rssFeed)

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(s.state.Current_user_name) == 0 {
		return errors.New("login is required.")
	}

	user, err := s.db.GetUser(context.Background(), s.state.Current_user_name)

	if err != nil {
		return fmt.Errorf("error getting the user id: %s", err)
	}

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

func handlerFollowFeed(s *state, cmd command) error {
	if len(s.state.Current_user_name) == 0 {
		return errors.New("login is required.")
	}

	user, err := s.db.GetUser(context.Background(), s.state.Current_user_name)

	if err != nil {
		return fmt.Errorf("error getting the user id: %s", err)
	}

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

func handlerFetchFollowing(s *state, cmd command) error {
	if len(s.state.Current_user_name) == 0 {
		return errors.New("login is required.")
	}

	user, err := s.db.GetUser(context.Background(), s.state.Current_user_name)

	if err != nil {
		return fmt.Errorf("error getting the user id: %s", err)
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)

	if err != nil {
		return fmt.Errorf("error getting the following feeds: %s", err)
	}

	fmt.Println(feeds)

	return nil
}


