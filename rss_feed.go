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

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request:%w", err)
	}
	req.Header.Add("User-Agent", "gator")

	// Create a new client and send a request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("faile to get response from client:%w", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error in getting data from response body:%w", err)
	}
	rssf := RSSFeed{}
	err = xml.Unmarshal(data, &rssf)
	if err != nil {
		return nil, fmt.Errorf("error in unmarshal data into struct")
	}

	rssf.Channel.Title = html.UnescapeString(rssf.Channel.Title)
	rssf.Channel.Description = html.UnescapeString(rssf.Channel.Description)

	for i := range rssf.Channel.Item {
		rssf.Channel.Item[i].Title = html.UnescapeString(rssf.Channel.Item[i].Title)
		rssf.Channel.Item[i].Title = html.UnescapeString(rssf.Channel.Item[i].Title)
	}

	return &rssf, nil

}
