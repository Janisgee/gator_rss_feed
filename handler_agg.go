package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()
	feedURL := "https://www.wagslane.dev/index.xml"

	rssf, err := fetchFeed(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("error in fetching the feed at %s: %w", feedURL, err)
	}
	fmt.Printf("Result struct: %v\n", rssf)

	return nil
}
