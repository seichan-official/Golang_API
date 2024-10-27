
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/joho/godotenv"
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
    // YouTube APIクライアントを作成
    client := &http.Client{
        Transport: &transport.APIKey{Key: apiKey},
    }
    service, err := youtube.New(client)
    if err != nil {
        return nil, fmt.Errorf("Error creating new YouTube client: %v", err)
    }

    // 検索クエリを設定
    call := service.Search.List([]string{"snippet"}).
        Q(query).   // 検索クエリ
        MaxResults(1). // 検索結果の最大数を1件に
        Type("video")  // 動画のみに絞る

    // APIリクエストを実行
    response, err := call.Do()
    if err != nil {
        return nil, fmt.Errorf("Error making YouTube API call: %v", err)
    }

    // 検索結果を格納するスライス
    var results []SearchResult

    // 検索結果を整形して追加
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

func handleSearch(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // POSTボディからクエリを取得
    var searchRequest SearchRequest
    err := json.NewDecoder(r.Body).Decode(&searchRequest)
    if err != nil || searchRequest.Query == "" {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // .envファイルの読み込み
    err = godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    // 環境変数からAPIキーを取得
    apiKey := os.Getenv("youtube_api_key")
    if apiKey == "" {
        log.Fatalf("API key not found in environment variables")
    }

    // YouTube APIを使って検索を実行
    results, err := searchYouTube(searchRequest.Query, apiKey)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 検索結果をJSON形式で返す
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(results)
}

func main() {
    // エンドポイントの登録
    http.HandleFunc("/api/youtube/search", handleSearch)

    // サーバーの起動
    fmt.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
