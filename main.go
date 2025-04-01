package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

// Apicurio Schema Registry URL
const schemaRegistryURL = "http://localhost:9090/api/artifacts/my-schema"

// apiConfig struct stores the database connection instance
type apiConfig struct {
	DB *database.Queries
}

// Fetch schema from Apicurio Registry
func fetchSchema() (map[string]interface{}, error) {
	resp, err := http.Get(schemaRegistryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var schema map[string]interface{}
	err = json.Unmarshal(body, &schema)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

// Validate incoming JSON data against schema
func validateJSON(data []byte, schema map[string]interface{}) error {
	var jsonData map[string]interface{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return fmt.Errorf("invalid JSON format")
	}

	// Example: Check if 'title' field exists
	if _, exists := jsonData["title"]; !exists {
		return fmt.Errorf("missing required field: title")
	}

	return nil
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read required environment variables
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("$PORT must be set")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	// Establish database connection
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	// Initialize database queries
	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}

	// Start background scraping
	go startScrapping(db, 10, time.Minute)

	// Initialize router
	router := chi.NewRouter()

	// Enable CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	// Define API routes
	v1Router := chi.NewRouter()

	// Health check
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handleErr)

	// User management
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	// Feed management
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	// Fetching posts for user
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	// Feed follow/unfollow
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollows))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	// Schema validation endpoint
	v1Router.Post("/validate", func(w http.ResponseWriter, r *http.Request) {
		schema, err := fetchSchema()
		if err != nil {
			http.Error(w, "Failed to fetch schema", http.StatusInternalServerError)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		err = validateJSON(body, schema)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte("JSON is valid!"))
	})

	// Mount API routes
	router.Mount("/v1", v1Router)

	// Start server
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Listening on port %s\n", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
