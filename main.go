package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Define your YouTube API key here (replace with your actual key)
	apiKey := os.Getenv("API_KEY")
	// Search query for football videos
	searchQuery := "football"
	// Maximum number of results to retrieve (limited to 10)
	maxResults := 10

	// Create a YouTube service object
	youtubeService, err := youtube.NewService(context.Background(), option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	// Define the search parameters
	call := youtubeService.Search.List([]string{"id", "snippet"}).Q(searchQuery).MaxResults(int64(maxResults)).Type("video").Order("date")

	// Execute the search request
	response, err := call.Do()
	if err != nil {
		log.Fatal(err)
	}

	// Print information about the first 5 videos
	if len(response.Items) > 0 {
		fmt.Println("Found", len(response.Items), "videos:")
		for _, item := range response.Items[:maxResults] {
			if snippet := item.Snippet; snippet != nil {
				fmt.Println("  - Title:", snippet.Title)
				fmt.Println("  - Description:", snippet.Title)
				fmt.Println("  - PublishedAt:", snippet.Title)
				fmt.Println("  - Thumbnail Url", snippet.Thumbnails.Default.Url)
				fmt.Println("  - ID:", item.Id)
			}
		}
	} else {
		fmt.Println("No videos found for", searchQuery)
	}
}
