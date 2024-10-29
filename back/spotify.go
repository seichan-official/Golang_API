package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
)


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
        http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
        return
    }

    spotifyToken = token // トークンを変数に保存

    response := map[string]string{
        "message": "認証に成功しました。",
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// Spotifyからユーザープロファイルを取得
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
        profile["image_url"] = firstImage["url"].(string)
    } else {
        profile["image_url"] = ""
    }

    return profile, nil
}

// ユーザーの最近の再生履歴を取得
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
        trackData := map[string]interface{}{
            "track_name": track["name"].(string),
            "album": map[string]interface{}{
                "name": track["album"].(map[string]interface{})["name"].(string),
                "image_url": track["album"].(map[string]interface{})["images"].([]interface{})[0].(map[string]interface{})["url"].(string),
            },
        }
        tracks = append(tracks, trackData)
    }

    return tracks, nil
}

// トークンと一緒にユーザーの再生履歴を取得するエンドポイント
func HandleUserHistory(w http.ResponseWriter, r *http.Request) {
    if spotifyToken == nil {
        http.Redirect(w, r, "/api/spotify/login", http.StatusTemporaryRedirect)
        return
    }

    client := spotifyConfig.Client(context.Background(), spotifyToken)
    profile, err := GetUserProfile(client)
    if err != nil {
        log.Printf("Error fetching user profile: %v\n", err)
        http.Error(w, "Error fetching profile", http.StatusInternalServerError)
        return
    }

    recentTracks, err := GetUserRecentPlayed(client)
    if err != nil {
        log.Printf("Error fetching recent tracks: %v\n", err)
        http.Error(w, "Error fetching recent tracks", http.StatusInternalServerError)
        return
    }

    data := map[string]interface{}{
        "user_profile":       profile,
        "recently_played_tracks": recentTracks,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}
