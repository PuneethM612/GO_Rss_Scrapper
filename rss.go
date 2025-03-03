package main

import (
	"encoding/xml" // Used for parsing XML data
	"io"           // Provides utilities for reading data
	"net/http"     // Handles HTTP requests
	"time"         // Used for setting timeouts
)

// RSSFeed struct represents the structure of an RSS feed.
type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`       // Title of the feed
		Link        string    `xml:"link"`        // Link to the website
		Description string    `xml:"description"` // Description of the feed
		Language    string    `xml:"language"`    // Language of the feed
		Items       []RSSItem `xml:"item"`        // List of RSS items (posts)
	} `xml:"channel"` // XML tag that matches the RSS structure
}

// RSSItem struct represents an individual item (post) in an RSS feed.
type RSSItem struct {
	Title       string `xml:"title"`       // Title of the post
	Link        string `xml:"link"`        // URL of the post
	Description string `xml:"description"` // Short summary of the post
	PubDate     string `xml:"pubDate"`     // Published date in string format
}

// urlTofeed fetches and parses an RSS feed from a given URL.
func urlTofeed(url string) (RSSFeed, error) {
	// Create an HTTP client with a timeout of 10 seconds and if the server takes more time to respond more than 10 seconds, it will timeout
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Send a GET request to the RSS feed URL
	resp, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err // Return an empty RSSFeed and the error if the request fails
	}
	defer resp.Body.Close() // Ensure the response body is closed after function execution

	// Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSFeed{}, err // Return an error if reading fails
	}

	// Initialize an empty RSSFeed struct
	rssFeed := RSSFeed{}

	// Unmarshal (convert) XML data into the RSSFeed struct
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return RSSFeed{}, err // Return an error if parsing XML fails
	}

	// Return the parsed RSS feed
	return rssFeed, nil
}
