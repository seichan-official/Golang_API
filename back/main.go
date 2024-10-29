package main

import (
    "os"
    "log"
    "net/http"

    "github.com/joho/godotenv"
    "github.com/rs/cors"
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

    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"}, // フロントエンドのURL
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        AllowCredentials: false,
    })

    mux := http.NewServeMux()
    // Spotify APIエンドポイントの設定
    mux.HandleFunc("/api/spotify/login", handleLogin)
    mux.HandleFunc("/callback", handleCallback)
    mux.HandleFunc("/api/spotify/history", HandleUserHistory)
    // YouTube APIエンドポイントの設定
    mux.HandleFunc("/api/youtube/search", YouTubeSearchHandler)



    handler := c.Handler(mux)

    // サーバーの起動
    log.Printf("Server started at http://localhost%s/", port)
    if err := http.ListenAndServe(port, handler); err != nil {
        log.Fatal(err)
    }
}
