package video

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lakshya8066/youtube_api/elasticsearch"
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

	err = StoreVideo(response, maxResults, db)
	if err != nil {
		fmt.Println(err)
	}

	elasticsearch.StoreToElastic()
	return nil
}
