package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Janisgee/gator_rss_feed/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	ctx := context.Background()

	limitPost := int32(2)
	// Set limits
	if len(cmd.Args) == 1 {
		num, err := strconv.ParseInt(cmd.Args[0], 10, 32)
		if err != nil {
			return fmt.Errorf("error in converting argument into number. %w", err)
		}
		limitPost = int32(num)
	}

	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limitPost),
	}
	// Get posts table
	posts, err := s.db.GetPostsForUser(ctx, params)
	if err != nil {
		return fmt.Errorf("error in getting posts:%w", err)
	}

	fmt.Printf("Total found %d posts, listed limit at %d posts.\n", len(posts), limitPost)
	for i := range posts {
		fmt.Printf("Title:           %s\n", posts[i].Title)
		fmt.Printf("Link:            %s \n", posts[i].Url)
		fmt.Printf("Publisted Date:  %v\n", posts[i].PublishedAt)
		fmt.Println()
		fmt.Println("=============================================")
	}

	return nil
}
