package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "github.com/google/uuid"
)

type UserProfile struct {
    ID              string `json:"id"`
    DisplayName     string `json:"display_name"`
    ProfileImageURL string `json:"profile_image_url"`
    Country         string `json:"country"`
}

type Artist struct {
    Name             string `json:"name"`
    SpotifyURL       string `json:"spotify_url"`
    SmallestImageURL string `json:"smallest_image_url"`
}

type Album struct {
    Name             string `json:"name"`
    SmallestImageURL string `json:"smallest_image_url"`
}

type Track struct {
    TrackName          string   `json:"track_name"`
    Album              Album    `json:"album"`
    Artists            []Artist `json:"artists"`
    YouTubeSearchQuery string   `json:"youtube_search_query"`
}

type UserHistoryResponse struct {
    Profile              UserProfile `json:"profile"`
    RecentlyPlayedTracks []Track     `json:"recently_played_tracks"`
}

// Spotify 認証ページへのリダイレクト
func handleLogin(w http.ResponseWriter, r *http.Request) {
    url := spotifyConfig.AuthCodeURL(oauthStateString)
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// 認証後のコールバック処理
func handleCallback(w http.ResponseWriter, r *http.Request) {
    state := r.FormValue("state")
    if state != oauthStateString {
        log.Println("invalid oauth state")
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }

    code := r.FormValue("code")
    token, err := spotifyConfig.Exchange(context.Background(), code)
    if err != nil {
        log.Printf("Code exchange failed: %v\n", err)
        // http.Redirect(w, r, "http://localhost:3000/main", http.StatusTemporaryRedirect)
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }

    sessionID := uuid.NewString() // 新しいセッションIDを生成
    sessionTokens[sessionID] = token // トークンをセッションに保存
    log.Printf("SessionID: %v\n\n", sessionID)
    log.Printf("token: %v\n\n", token)
    http.SetCookie(w, &http.Cookie{
        Name:  "session_id",
        Value: sessionID,
        Path:  "/",
    })

    response := map[string]string{
        "message": "認証に成功しました。",
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// Spotifyのユーザープロファイルを取得する関数
func GetUserProfile(client *http.Client) (map[string]string, error) {
    resp, err := client.Get("https://api.spotify.com/v1/me")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    profile := map[string]string{
        "id":           result["id"].(string),
        "display_name": result["display_name"].(string),
        "country":      result["country"].(string),
    }

    if images, ok := result["images"].([]interface{}); ok && len(images) > 0 {
        firstImage := images[0].(map[string]interface{})
        profile["profile_image_url"] = firstImage["url"].(string)
    } else {
        profile["profile_image_url"] = ""
    }

    return profile, nil
}



func HandleUserHistory(w http.ResponseWriter, r *http.Request) {

    
    // セッションIDを取得
    cookie, err := r.Cookie("session_id")
    if err != nil || sessionTokens[cookie.Value] == nil {
        http.Redirect(w, r, "/api/spotify/login", http.StatusTemporaryRedirect)
        return
    }
    log.Printf("SessionTokens: %v\n\n", sessionTokens)



    spotifyToken := sessionTokens[cookie.Value]
    log.Printf("Cookieから取得したSpotifyToken: %v\n\n", spotifyToken)
    

    if spotifyToken == nil {
        http.Redirect(w, r, "/api/spotify/login", http.StatusTemporaryRedirect)
        return
    }

    client := spotifyConfig.Client(context.Background(), spotifyToken)
    profileData, err := GetUserProfile(client)
    if err != nil {
        log.Printf("Error fetching user profile: %v\n", err)
        http.Error(w, "Error fetching profile", http.StatusInternalServerError)
        return
    }

    recentTracksData, err := GetUserRecentPlayed(client)
    if err != nil {
        log.Printf("Error fetching recent tracks: %v\n", err)
        http.Error(w, "Error fetching recent tracks", http.StatusInternalServerError)
        return
    }

    // プロフィールデータの組み立て
    userProfile := UserProfile{
        ID:              profileData["id"],
        DisplayName:     profileData["display_name"],
        ProfileImageURL: profileData["image_url"],
        Country:         profileData["country"],
    }

    // トラックデータの組み立て
    var tracks []Track
    for _, item := range recentTracksData {
        trackName := item["track_name"].(string)

        // アルバム情報の取得
        albumData := item["album"].(map[string]interface{})
        album := Album{
            Name:             albumData["name"].(string),
            SmallestImageURL: albumData["smallest_image_url"].(string),
        }

        // アーティスト情報の取得
        var artists []Artist
        for _, artistItem := range item["artists"].([]map[string]string) {
            artist := Artist{
                Name:             artistItem["name"],
                SpotifyURL:       artistItem["spotify_url"],
                SmallestImageURL: artistItem["smallest_image_url"],
            }
            artists = append(artists, artist)
        }

        youtubeSearchQuery := item["youtube_search_query"].(string)

        track := Track{
            TrackName:          trackName,
            Album:              album,
            Artists:            artists,
            YouTubeSearchQuery: youtubeSearchQuery,
        }
        tracks = append(tracks, track)
    }

    // レスポンスデータの組み立て
    response := UserHistoryResponse{
        Profile:              userProfile,
        RecentlyPlayedTracks: tracks,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}


func GetUserRecentPlayed(client *http.Client) ([]map[string]interface{}, error) {
    resp, err := client.Get("https://api.spotify.com/v1/me/player/recently-played")
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    var tracks []map[string]interface{}
    for _, item := range result["items"].([]interface{}) {
        track := item.(map[string]interface{})["track"].(map[string]interface{})
        album := track["album"].(map[string]interface{})
        artists := track["artists"].([]interface{})

        // アーティスト情報を取得
        var artistData []map[string]string
        for _, artist := range artists {
            artistInfo := artist.(map[string]interface{})
            artistData = append(artistData, map[string]string{
                "name":             artistInfo["name"].(string),
                "spotify_url":      fmt.Sprintf("https://open.spotify.com/artist/%s", artistInfo["id"].(string)),
                "smallest_image_url": "", // 必要なら取得
            })
        }

        // アルバムの画像URLのうち一番小さいものを取得
        imageURL := ""
        if images, ok := album["images"].([]interface{}); ok && len(images) > 0 {
            lastImage := images[len(images)-1].(map[string]interface{})
            imageURL = lastImage["url"].(string)
        }

        // YouTube検索用クエリを生成
        youtubeQuery := fmt.Sprintf("%s [%s] Official Music Video", track["name"].(string), artistData[0]["name"])

        trackData := map[string]interface{}{
            "track_name": track["name"].(string),
            "album": map[string]interface{}{
                "name":              album["name"].(string),
                "smallest_image_url": imageURL,
            },
            "artists":              artistData,
            "youtube_search_query": youtubeQuery,
        }
        tracks = append(tracks, trackData)
    }

    return tracks, nil
}
