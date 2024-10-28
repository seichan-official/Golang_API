package main

import (
    "os"
    "log"
    "net/http"

    "github.com/joho/godotenv"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/spotify"
)


var spotifyConfig *oauth2.Config  // グローバル変数として宣言
var oauthStateString = "random"   // 任意の状態文字列
var spotifyToken *oauth2.Token    // Spotify認証トークンを保持する変数

func main() {
    if err := godotenv.Load("./.env"); err != nil {
        log.Println("Error loading .env file:", err)
    }

    port := ":8080"


    // godotenv.Load()の後でspotifyConfigを初期化
    spotifyConfig = &oauth2.Config{
        ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
        ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
        RedirectURL:  "http://localhost:8080/callback",
        Endpoint:     spotify.Endpoint,
        Scopes: []string{
            "user-read-private", "user-read-email",
            "user-read-recently-played",
        },
    }

    // Spotify APIエンドポイントの設定
    http.HandleFunc("/api/spotify/login", handleLogin)
    http.HandleFunc("/callback", handleCallback)
    http.HandleFunc("/api/spotify/history", HandleUserHistory)
    // YouTube APIエンドポイントの設定
    http.HandleFunc("/api/youtube/search", YouTubeSearchHandler)

    // サーバーの起動
    log.Printf("Server started at http://localhost%s/", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatal(err)
    }
}
