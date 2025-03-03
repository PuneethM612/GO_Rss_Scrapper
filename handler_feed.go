package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/PuneethM06/rssagg/internal/database"
	"github.com/google/uuid"
)

// handlerCreateFeed handles creating a new feed.
func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	// Define expected JSON request structure
	type parameters struct {
		Name string `json:"name"` // Feed name
		URL  string `json:"url"`  // Feed URL
	}

	// Decode request body into parameters struct
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		// If decoding fails, return an error response
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Create a new feed entry in the database
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),                                // Generate unique feed ID
		CreatedAt: time.Now().UTC(),                          // Set creation timestamp
		UpdatedAt: time.Now().UTC(),                          // Set update timestamp
		Name:      params.Name,                               // Store feed name
		Url:       params.URL,                                // Store feed URL
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true}, // Store user ID (nullable)
	})

	if err != nil {
		log.Printf("Error creating feed: %v", err) // Log error
		respondWithError(w, http.StatusInternalServerError, "Unable to create feed")
		return
	}

	// Send the created feed as a response
	respondwithJSON(w, http.StatusOK, feed)
}

// handlerGetFeeds retrieves all feeds from the database.
func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	// Fetch all feeds from the database
	feed, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		log.Printf("Error getting feeds: %v", err) // Log error
		respondWithError(w, http.StatusInternalServerError, "Unable to get feeds")
		return
	}

	// Send the retrieved feeds as a JSON response
	respondwithJSON(w, http.StatusOK, databaseFeedsToFeeds(feed))
}
