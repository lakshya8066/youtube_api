package video

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/youtube/v3"
)

func createTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS videos (
        id INT AUTO_INCREMENT PRIMARY KEY,
		video_id VARCHAR(50),
        title VARCHAR(255),
        description TEXT,
        published_at DATETIME,
        thumbnail_url VARCHAR(255)
    )`)
	return err
}

func StoreVideo(response *youtube.SearchListResponse, maxResults int, db *sql.DB) error {
	// Create the videos table if it doesn't exist
	err := createTable(db)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return err
	}

	fmt.Println("Saving videos")
	// Iterate through videos and store data in MySQL
	for _, item := range response.Items[:maxResults] {
		if snippet := item.Snippet; snippet != nil {
			publishedAt, err := time.Parse(time.RFC3339, snippet.PublishedAt)
			if err != nil {
				log.Println("Error parsing publishedAt:", err)
				continue // Skip this item if there's an error parsing publishedAt
			}
			// Prepare SQL statement with placeholders
			stmt, err := db.Prepare("INSERT INTO videos (video_id, title, description, published_at, thumbnail_url) VALUES (?, ?, ?, ?, ?)")
			if err != nil {
				fmt.Println(err)
				return err // Propagate error
			}
			defer stmt.Close()

			// Extract and insert video data
			_, err = stmt.Exec(item.Id.Kind+item.Id.VideoId, snippet.Title, snippet.Description, publishedAt, snippet.Thumbnails.Default.Url)
			if err != nil {
				fmt.Println(err)
				return err // Propagate error
			}
		}
	}
	fmt.Println("Successfully stored", len(response.Items), "videos in database.")
	return nil
}
