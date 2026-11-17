package main

func scrapeFeeds(s *state, cmd command, user database.User) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())

	if err != nil {
		return fmt.Errorf("error getting the next feed to fetch: %s", err)
	}

	payload := database.MarkFeedFetchedParams{
		LastFetchedAt: time.Now(),
		ID:            feed.ID,
	}

	if err := s.db.MarkFeedFetched(context.Background(), payload); err != nil {
		return fmt.Errorf("error marking the feed: %s", err)
	}

	rssFeed, err := fetchFeed(context.Background(), feed.url)

	for _, item := range rssFeed.Item {
		fmt.Println(item.Title)
	}

	return nil
}
