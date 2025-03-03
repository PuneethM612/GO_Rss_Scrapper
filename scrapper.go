package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/PuneethM06/rssagg/internal/database" // Importing the database package
	"github.com/google/uuid"                         // Importing UUID for generating unique IDs
)

// startScrapping initiates the RSS feed scraping process at regular intervals.
func startScrapping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)

	// Ticker triggers the scraping process at the specified time interval
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C { // Infinite loop to keep running the scraper
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency), // Fetch a batch of feeds based on concurrency level
		)
		if err != nil {
			log.Println(err)
			continue // Skip this iteration if there's an error
		}

		// Waitgroup ensures in hadnling multiple tasks and manage multiple goroutines in batch systems.
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)                   // Increase counter for each feed being processed
			go scrapeFeed(db, wg, feed) // Launch a goroutine to scrape the feed concurrently
		}
		wg.Wait() // Wait until all goroutines complete before proceeding
	}
}

// scrapeFeed fetches and processes an individual RSS feed.
func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done() // Decrement the WaitGroup counter when the function completes

	// Validate that the feed URL is not empty
	if feed.Url == "" {
		log.Println("Skipping feed: Empty URL")
		return
	}

	// Mark the feed as fetched in the database to prevent duplicate processing
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	} // Parses in terms of it converts the XML file into the structres that we can understand

	// Fetch and parse the RSS feed
	rssFeed, err := urlTofeed(feed.Url)
	if err != nil {
		log.Println("Error parsing feed:", err)
		return
	}

	for _, item := range rssFeed.Channel.Items { // Iterate over all items (posts) in the feed
		// Handle cases where the description might be null
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		// Parse the publication date of the post
		var pubAt time.Time
		dateFormats := []string{
			time.RFC1123Z, time.RFC1123, time.RFC3339, "Mon, 2 Jan 2006 15:04:05 MST",
		}

		// Try parsing the date using multiple formats
		for _, format := range dateFormats {
			pubAt, err = time.Parse(format, item.PubDate)
			if err == nil {
				break // Exit loop if parsing is successful
			}
		}

		// If parsing fails, default to the current time
		if err != nil {
			log.Println("Error parsing pub date:", err, "Raw Date:", item.PubDate)
			pubAt = time.Now().UTC()

		}

		// Insert the parsed feed item into the database
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(), // Generate a unique ID for the post
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID, // Associating the post with its feed
		})
		if err != nil {
			// Skip duplicate entries to avoid inserting the same post multiple times
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("Error creating post:", err)
		}
	}
}
