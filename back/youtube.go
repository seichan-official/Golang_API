package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "google.golang.org/api/googleapi/transport"
    "google.golang.org/api/youtube/v3"
)

// 検索リクエスト用の構造体
type SearchRequest struct {
    Query string `json:"query"`
}

// 検索結果を返すための構造体
type SearchResult struct {
    Title       string `json:"title"`
    VideoURL    string `json:"video_url"`
    Thumbnail   string `json:"thumbnail"`
    Description string `json:"description"`
    Channel     string `json:"channel"`
}

func searchYouTube(query string, apiKey string) ([]SearchResult, error) {
    client := &http.Client{
        Transport: &transport.APIKey{Key: apiKey},
    }
    service, err := youtube.New(client)
    if err != nil {
        return nil, fmt.Errorf("Error creating new YouTube client: %v", err)
    }

    call := service.Search.List([]string{"snippet"}).
        Q(query).
        MaxResults(1).
        Type("video")

    response, err := call.Do()
    if err != nil {
        return nil, fmt.Errorf("Error making YouTube API call: %v", err)
    }

    var results []SearchResult
    for _, item := range response.Items {
        videoID := item.Id.VideoId
        videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
        result := SearchResult{
            Title:       item.Snippet.Title,
            VideoURL:    videoURL,
            Thumbnail:   item.Snippet.Thumbnails.High.Url,
            Description: item.Snippet.Description,
            Channel:     item.Snippet.ChannelTitle,
        }
        results = append(results, result)
    }

    return results, nil
}

func YouTubeSearchHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    var searchRequest SearchRequest
    err := json.NewDecoder(r.Body).Decode(&searchRequest)
    if err != nil || searchRequest.Query == "" {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    apiKey := os.Getenv("youtube_api_key")
    if apiKey == "" {
        http.Error(w, "API key not found in environment variables", http.StatusInternalServerError)
        return
    }

    results, err := searchYouTube(searchRequest.Query, apiKey)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}
