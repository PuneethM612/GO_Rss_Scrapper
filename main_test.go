package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// TestRouterSetup verifies the router is set up correctly
func TestRouterSetup(t *testing.T) {
	r := chi.NewRouter()

	// Apply CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	// Sample route
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	// Create a test server
	req, _ := http.NewRequest("GET", "/ping", nil)
	recorder := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(recorder, req)

	// Validate response
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", recorder.Code)
	}

	if recorder.Body.String() != "pong" {
		t.Errorf("Expected response 'pong', got '%s'", recorder.Body.String())
	}
}
