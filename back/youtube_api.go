package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
    "google.golang.org/api/googleapi/transport"
    "google.golang.org/api/youtube/v3"
)

func main() {
    // .envファイルの読み込み
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    // 環境変数からAPIキーを取得
    apiKey := os.Getenv("youtube_api_key")
    if apiKey == "" {
        log.Fatalf("API key not found in environment variables")
    }

    // YouTubeサービスを作成
    client := &http.Client{
        Transport: &transport.APIKey{Key: apiKey},
    }
    service, err := youtube.New(client)
    if err != nil {
        log.Fatalf("Error creating new YouTube client: %v", err)
    }

    // 検索クエリを設定
    query := "あいみょん マリーゴールド Music Video"
    call := service.Search.List([]string{"snippet"}).
        Q(query).   // 検索クエリ
        MaxResults(5). // 検索結果の最大数を5件に
        Type("video") // 動画のみに絞る

    // APIリクエストを実行
    response, err := call.Do()
    if err != nil {
        log.Fatalf("Error making YouTube API call: %v", err)
    }

    // 検索結果を表示
    fmt.Println("Search results:")
    for _, item := range response.Items {
        title := item.Snippet.Title
        videoID := item.Id.VideoId
        thumbnailURL := item.Snippet.Thumbnails.High.Url // 高解像度サムネイルを取得
        description := item.Snippet.Description
        channelTitle := item.Snippet.ChannelTitle
        videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)

        // 結果を整形して表示
        fmt.Printf("Title: %s\nVideo URL: %s\nThumbnail: %s\nDescription: %s\nChannel: %s\n\n",
            title, videoURL, thumbnailURL, description, channelTitle)
    }
}
