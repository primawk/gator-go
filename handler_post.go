package main

import (
	"context"
	"fmt"
	"strconv"
)

func handlerFetchPosts(s *state, cmd command) error {
	var limit int32 = 2

	if len(cmd.args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("invalid limit %q: %w", cmd.args[0], err)
		}
		limit = int32(parsedLimit)
	}

	posts, err := s.db.GetPosts(context.Background(), limit)

	if err != nil {
		return fmt.Errorf("error fetching posts: %s\n", err)
	}

	for _, post := range posts {
		fmt.Printf("- %s\n", post)
	}

	

	return nil
}
