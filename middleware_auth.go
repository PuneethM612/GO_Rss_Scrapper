package main

import (
	"net/http"

	"github.com/PuneethM06/rssagg/internal/auth"     // Import authentication helper functions
	"github.com/PuneethM06/rssagg/internal/database" // Import database operations
)

// Define a custom handler type that includes the authenticated user
type authHandler func(http.ResponseWriter, *http.Request, database.User)

// Middleware function to authenticate requests using an API key
func (apiCfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header) // It is basically calling the API key for authorizing
		if err != nil {
			// If the API is not authorized, return an authorized error.
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// Retrieve the data from the database using the API key
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			// If user retrieval fails, return an internal server error
			respondWithError(w, http.StatusInternalServerError, "Unable to get user")
			return
		}

		// call the actual handler with the authenticated user
		handler(w, r, user)
	}
}
