package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// respondWithError sends a JSON-formatted error response with a given status code.
func respondWithError(w http.ResponseWriter, code int, msg string) {
	// Log server-side errors (5xx) for debugging
	if code > 499 { // Client-side errors (4xx) are usually not logged, but server errors (5xx) are.
		log.Println("Responding with 5xx error:", msg)
	}

	// Define the error response structure
	type errorResponse struct {
		Error string `json:"error"` // JSON key for error message
	}

	// Send the error response as JSON
	respondwithJSON(w, code, errorResponse{Error: msg})
}

// respondwithJSON sends a JSON response with the given payload and status code.
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Convert the payload into JSON format
	// In this case the payload is been refered to the data that needs to be converted into the JSON format.
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Panicln("Unable to marshal JSON") // Log critical error and stop execution
		return
	}

	w.Header().Add("Content-Type", "application/json") // Indicates the user that the data is in json format.

	// Set the HTTP status code
	w.WriteHeader(code)

	// Write the JSON response data to the client
	w.Write(dat)
}
