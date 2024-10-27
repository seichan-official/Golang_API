package main

import (
    "log"
    "net/http"
)

func main() {
    port := ":8080"

    // 静的ファイルサーバーとして`/docs`を提供
    fs := http.FileServer(http.Dir("./docs"))
    http.Handle("/docs/", http.StripPrefix("/docs/", fs))  // URLパスの/docs/部分を除去

    // Spotify APIエンドポイントの設定
    http.HandleFunc("/api/spotify/login", SpotifyLoginHandler)
    http.HandleFunc("/api/spotify/history", SpotifyHistoryHandler)

    // YouTube APIエンドポイントの設定
    http.HandleFunc("/api/youtube/search", YouTubeSearchHandler)

    // サーバーの起動
    log.Printf("Server started at http://localhost%s/", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        log.Fatal(err)
    }
}
