package elasticsearch

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	_ "github.com/go-sql-driver/mysql"
)

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

func StoreToElastic() {
	// Connect to MySQL database
	db, err := dbConnection()
	if err != nil {
		log.Fatalf("Error connecting to MySQL: %v", err)
	}
	defer db.Close()

	// Query videos from MySQL
	rows, err := db.Query("SELECT * FROM videos")
	if err != nil {
		log.Fatalf("Error querying videos from MySQL: %v", err)
	}
	defer rows.Close()

	// Connect to Elasticsearch
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	// Iterate over rows and load data into Elasticsearch
	for rows.Next() {
		var video Video
		if err := rows.Scan(&video.ID, &video.VideoID, &video.Title, &video.Description, &video.PublishedAt, &video.ThumbnailURL); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		// Transform data if necessary
		// For simplicity, let's assume no transformation is needed

		// Index data into Elasticsearch
		if err := indexDocument(es, video); err != nil {
			log.Printf("Error indexing document into Elasticsearch: %v", err)
			continue
		}
	}
}

func indexDocument(es *elasticsearch.Client, video Video) error {
	// Marshal video data to JSON
	data, err := json.Marshal(video)
	if err != nil {
		return fmt.Errorf("error marshaling video data: %w", err)
	}

	// Index document into Elasticsearch
	_, err = es.Index("videos", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("error indexing document into Elasticsearch: %w", err)
	}

	return nil
}
