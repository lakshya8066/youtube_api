package video

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/api/youtube/v3"
)

func FetchVideos(youtubeService *youtube.Service, searchQuery string, maxResults int, db *sql.DB) error {
	// Define the search parameters
	call := youtubeService.Search.List([]string{"id", "snippet"}).Q(searchQuery).MaxResults(int64(maxResults)).Type("video").Order("date")

	// Execute the search request
	response, err := call.Do()
	if err != nil {
		fmt.Println(err)
	}

	StoreVideo(response, maxResults, db)

	// elasticsearch.StoreToElastic()

	// Print information about the first 5 videos
	// if len(response.Items) > 0 {
	// 	fmt.Println("Found", len(response.Items), "videos:")
	// 	for _, item := range response.Items[:maxResults] {
	// 		if snippet := item.Snippet; snippet != nil {
	// 			fmt.Println("  - Title:", snippet.Title)
	// 			fmt.Println("  - Description:", snippet.Title)
	// 			fmt.Println("  - PublishedAt:", snippet.Title)
	// 			fmt.Println("  - Thumbnail Url", snippet.Thumbnails.Default.Url)
	// 			fmt.Println("  - ID:", item.Id)
	// 		}
	// 	}
	// } else {
	// 	fmt.Println("No videos found for", searchQuery)
	// }
	return nil
}
