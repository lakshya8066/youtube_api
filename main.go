package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lakshya8066/youtube_api/controllers"
	"github.com/lakshya8066/youtube_api/video"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func fetchVideos(youtubeService *youtube.Service, searchQuery string, maxResults int, db *sql.DB) {
	for {
		// Code to be executed every 10 seconds
		fmt.Println("This message is printed every 10 seconds.")

		video.FetchVideos(youtubeService, searchQuery, maxResults, db)

		// Sleep for 10 seconds using time.Sleep
		time.Sleep(10 * time.Second)
	}
}

func dbConnection() (*sql.DB, error) {
	username := "root"
	password := "12345678"
	host := "localhost"
	port := "3307"
	database := "youtube"

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

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Define your YouTube API key here (replace with your actual key)
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	// Search query for football videos
	searchQuery := "football"
	// Maximum number of results to retrieve (limited to 10)
	maxResults := 10

	// Create a YouTube service object
	youtubeService, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	db, err := dbConnection()

	go fetchVideos(youtubeService, searchQuery, maxResults, db)

	defer db.Close()

	controllers.Handler()
}
