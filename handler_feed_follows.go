package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/PuneethM06/rssagg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler to create a new feed follow for a user
func (apiCfg *apiConfig) handlerCreateFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"` // Feed ID provided in JSON request body
	}

	// Decode JSON request body into parameters struct
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Create a new feed follow entry in the database
	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),       // Generate a new UUID for feed follow entry
		CreatedAt: time.Now().UTC(), // Set creation timestamp
		UpdatedAt: time.Now().UTC(), // Set update timestamp
		UserID:    user.ID,          // Associate feed follow with the user
		FeedID:    params.FeedID,    // Associate feed follow with the provided Feed ID
	})

	if err != nil {
		log.Printf("Error creating feed follow: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to create feed follow")
		return
	}

	// Return success response with the created feed follow details
	respondwithJSON(w, http.StatusOK, databaseFeedFollowToFeedFollow(feedFollow))
}

// Handler to get all feed follows for a user
func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	// Fetch feed follows from the database for the given user
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil {
		log.Printf("Error getting feed follows: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to get feed follows")
		return
	}

	// Return success response with the list of feed follows
	respondwithJSON(w, http.StatusOK, databaseFeedFollowsToFeedFollows(feedFollows))
}

// Handler to delete a feed follow for a user
func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	// Extract feedFollowID from the URL
	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feedFollowID")
		return
	}

	// Delete the feed follow entry
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feedFollowID,
	})
	if err != nil {
		log.Printf("Error deleting feed follow: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to delete feed follow")
		return
	}

	// Check if there are any remaining followers for the feed
	count, err := apiCfg.DB.CountFeedFollows(r.Context(), feedFollowID)
	if err != nil {
		log.Printf("Error counting feed follows: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error checking remaining follows")
		return
	}

	// If no followers remain, delete the feed itself
	if count == 0 {
		err = apiCfg.DB.DeleteFeed(r.Context(), feedFollowID)
		if err != nil {
			log.Printf("Error deleting feed: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Unable to delete feed")
			return
		}
	}

	// Return success message
	respondwithJSON(w, http.StatusOK, map[string]string{"message": "Feed follow deleted successfully"})
}
