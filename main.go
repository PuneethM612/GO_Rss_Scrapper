package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuneethM06/rssagg/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// apiConfig struct stores the database connection instance
// This will be used to handle API requests that interact with the database
type apiConfig struct {
	DB *database.Queries
}

func main() {
	// Load the .env file to get environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read the PORT environment variable
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("$PORT must be set")
	}

	// Read the database connection URL from environment variables
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	// Establish a connection to the PostgreSQL database
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	// Initialize database queries
	db := database.New(conn)
	api_cfg := apiConfig{
		DB: db,
	}

	// Start the background process for scraping feeds periodically
	go startScrapping(db, 10, time.Minute)

	// Create a new router for handling HTTP requests
	router := chi.NewRouter()

	// Set up Cross-Origin Resource Sharing (CORS) middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	// Define a sub-router for API versioning (v1)
	v1Router := chi.NewRouter()

	// Health check endpoint
	v1Router.Get("/healthz", handlerReadiness)
	// Test error handling endpoint
	v1Router.Get("/err", handleErr)

	// User management endpoints
	v1Router.Post("/users", api_cfg.handlerCreateUser)                     // Create a new user
	v1Router.Get("/users", api_cfg.middlewareAuth(api_cfg.handlerGetUser)) // Get authenticated user details

	// Feed management endpoints
	v1Router.Post("/feeds", api_cfg.middlewareAuth(api_cfg.handlerCreateFeed)) // Create a feed
	v1Router.Get("/feeds", api_cfg.handlerGetFeeds)                            // Get all feeds

	// Fetching posts related to user's followed feeds
	v1Router.Get("/posts", api_cfg.middlewareAuth(api_cfg.handlerGetPostsForUser))

	// Feed follow/unfollow operations
	v1Router.Post("/feed_follows", api_cfg.middlewareAuth(api_cfg.handlerCreateFeedFollows))                 // Follow a feed
	v1Router.Get("/feed_follows", api_cfg.middlewareAuth(api_cfg.handlerGetFeedFollows))                     // Get followed feeds
	v1Router.Delete("/feed_follows/{feedFollowID}", api_cfg.middlewareAuth(api_cfg.handlerDeleteFeedFollow)) // Unfollow a feed

	// Mount v1Router under the /v1 path
	router.Mount("/v1", v1Router)

	// Create and start the HTTP server
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	// Log the port the server is listening on
	log.Printf("Listening on port %s\n", portString)

	// Start the server and log any errors
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Port: %s\n", portString)
}
