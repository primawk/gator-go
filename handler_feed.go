package main

import (
	"context"
	"encoding/xml"
	// "errors"
	"fmt"
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
