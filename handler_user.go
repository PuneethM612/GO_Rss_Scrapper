package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/PuneethM06/rssagg/internal/database" // Importing database package
	"github.com/google/uuid"                         // Importing UUID package to generate unique user IDs
)

// 🔹 Handler to Create a New User
func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// Define the expected JSON payload format
	type parameters struct {
		Name string `json:"name"` // Extract the "name" field from JSON
	}

	// Create a JSON decoder to parse the request body, and the data in r.body is from the request sent.
	decoder := json.NewDecoder(r.Body)

	// Initialize an empty parameters struct
	params := parameters{}
	err := decoder.Decode(&params) // storing the data decode in the decoder is stored in the params struct that is created.
	if err != nil {
		// Return an error response if JSON parsing fails
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Create a new user in the database
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),       // Generate a unique ID for the user
		CreatedAt: time.Now().UTC(), // Store creation timestamp
		UpdatedAt: time.Now().UTC(), // Store update timestamp
		Name:      params.Name,      // Store the provided user name
	})

	if err != nil {
		log.Printf("Error creating user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to create user")
		return
	}

	// Return a success response with the created user details
	respondwithJSON(w, http.StatusOK, databaseUserToUser(user))
}

// 🔹 Handler to Retrieve User Details
func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	// Convert the database user model to API response format and send it
	respondwithJSON(w, http.StatusOK, databaseUserToUser(user))
}

// 🔹 Handler to Retrieve Posts for a Specific User
func (apiCfg *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	// Fetch posts from the database for the given user
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{ // context is used for cancellng the database query in case it timeouts or user closes the request.
		UserID: user.ID, // Fetch posts only for the authenticated user
		Limit:  10,      // Limit the results to the latest 10 posts
	})

	if err != nil {
		log.Printf("Error getting posts for user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Unable to get posts for user")
		return
	}

	// Return the posts in JSON format
	respondwithJSON(w, http.StatusOK, databasePostsToPosts(posts))
}
