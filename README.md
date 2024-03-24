# youtube_api

The goal of this project is to create an API for fetching the latest videos from YouTube, sorted in reverse chronological order of their publishing date-time, based on a given tag or search query. The response is paginated to efficiently handle large amounts of data.

## Running the project

1. Clone the Project: Clone the project repository to your local machine.
```
git clone https://github.com/lakshya8066/youtube_api.git
```
2. Navigate to Project Directory: Move to the project directory.
```
cd youtube_api
```
3. Start Docker Containers: Run the following command to start MySQL, Elasticsearch, and Kibana containers using Docker Compose.
```
docker-compose up -d
```
This command initializes the necessary database and search engine components.
4. Set API Key: Place your YouTube API key in the .env file.
```
echo "YOUTUBE_API_KEY=<api_key>" > .env
```
5. Run the Project: Execute the Go application.
```
go run main.go
```

## API Endpoints

### 1. Fetch Videos:

- Endpoint: http://localhost:8080/videos?page={page_number}&page_size={page_size}
- Description: Returns a paginated response of videos stored in the MySQL database, sorted in descending order of published datetime.
- Example curl request:
```
curl --location 'http://localhost:8080/videos?page=2&page_size=10'
```

### 2. Update API Key:

- Endpoint: http://localhost:8080/api-key
- Description: Sends a POST request with the new API key when the old key's quota has been exhausted.
- Example curl request:

```curl
curl --location --request GET 'http://localhost:8080/api-key' \
--header 'Content-Type: application/json' \
--data '{
    "new_api_key":"anmol"
}'
```

## Architecture

- The application periodically fetches video data from the YouTube Data API v3 every 10 seconds.
- Fetched data is stored in the MySQL server and simultaneously sent to the Elasticsearch server for visualization in the Kibana Dashboard.
- Asynchronous Go routines are employed to handle data retrieval from YouTube, ensuring non-blocking operation.

## Optimizations

- Batch Processing: This can be implemented to store multiple videos in a single api call instead of fetching the data multiple times. This will reduce connection calls to the database.

- Removing MySQL: This app could also be build by directly indexing the video data into ElasticSearch as the requirements of the app is minimum. 