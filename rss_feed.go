package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
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

	// Create a new Client
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request:%w", err)
	}
	req.Header.Add("User-Agent", "gator")

	// Send a request
	resp, err := httpClient.Do(req)
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

	for i, item := range rssf.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssf.Channel.Item[i] = item
	}

	return &rssf, nil

}
