package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Video represents a video entry in the database
type Video struct {
	ID           int    `json:"id"`
	VideoID      string `json:"video_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	PublishedAt  string `json:"published_at"`
	ThumbnailURL string `json:"thumbnail_url"`
}

func dbConnection() (*sql.DB, error) {
	username := "root"
	password := "12345678"
	host := "localhost"
	port := "3306"
	database := "yt"

	// Form the connection string
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	// Open a connection to the MySQL database
	fmt.Println("Opening connection")
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("Opened connection")

	return db, nil
}

func getVideosHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	// Connect to the database
	db, err := dbConnection()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query the database for videos
	offset := (page - 1) * pageSize
	rows, err := db.Query("SELECT * FROM videos ORDER BY published_at DESC LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate over the rows and build the response
	var videos []Video
	for rows.Next() {
		var video Video
		if err := rows.Scan(&video.ID, &video.VideoID, &video.Title, &video.Description, &video.PublishedAt, &video.ThumbnailURL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		videos = append(videos, video)
	}

	// Encode the response as JSON and write it to the response writer
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

// APIKeyHandler handles requests to update the YouTube API key
func APIKeyHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get the new API key
	var requestData struct {
		NewAPIKey string `json:"new_api_key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// Update the stored API key with the new one
	mutex.Lock()
	defer mutex.Unlock()
	apiKey = requestData.NewAPIKey

	// Update .env file with the new API key
	if err := godotenv.Write(map[string]string{"YOUTUBE_API_KEY": apiKey}, ".env"); err != nil {
		log.Printf("Error writing to .env file: %s", err)
		http.Error(w, "Failed to update API key", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	response := struct {
		Message string `json:"message"`
	}{"YouTube API key updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	reloadEnv()
}

var (
	apiKey string
	mutex  sync.Mutex
)

func reloadEnv() {
	// Reload environment variables from .env file if it has been modified
	lastLoad := time.Now()
	if info, err := os.Stat(".env"); err == nil && info.ModTime().After(lastLoad) {
		if err := godotenv.Load(); err != nil {
			log.Printf("Error reloading .env file: %s", err)
		} else {
			apiKey = os.Getenv("YOUTUBE_API_KEY")
			lastLoad = time.Now()
		}
	}
}

func Handler() {
	router := mux.NewRouter()
	router.HandleFunc("/", helloHandler)
	router.HandleFunc("/videos", getVideosHandler).Methods("GET")
	router.HandleFunc("/api-key", APIKeyHandler).Methods("POST")

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}
