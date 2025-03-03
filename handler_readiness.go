package main

import "net/http"

// handlerReadiness checks if the server is ready to handle requests.
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	// Respond with HTTP 200 OK status and an empty JSON object.
	respondwithJSON(w, http.StatusOK, struct{}{})
}
