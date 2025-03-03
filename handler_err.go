package main

import "net/http"

// handleErr handles errors by responding with a generic internal server error message.
func handleErr(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Something went wrong", http.StatusInternalServerError)
}
