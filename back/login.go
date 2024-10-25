package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

var spotifyConfig *oauth2.Config
var token *oauth2.Token

func main() {
	// .envファイルの読み込み
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Spotify APIのOAuth2設定
	spotifyConfig = &oauth2.Config{
		ClientID:     os.Getenv("client_id"),   // Spotify Developerから取得
		ClientSecret: os.Getenv("secret_id"),   // Spotify Developerから取得
		RedirectURL:  "http://localhost:8080/callback", // リダイレクトURL
		Endpoint:     spotify.Endpoint,         // Spotify用のOAuth2エンドポイント
		Scopes: []string{
			"user-read-private",
			"user-read-email",
			"user-read-recently-played",
		},
	}

	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)
	http.HandleFunc("/profile", handleProfile)

	fmt.Println("Server started at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// /loginエンドポイント: Spotifyの認証ページにリダイレクト
func handleLogin(w http.ResponseWriter, r *http.Request) {
	oauthStateString := "test_state"

	// 認証URLを生成
	url := spotifyConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// /callbackエンドポイント: 認証後に呼び出される
func handleCallback(w http.ResponseWriter, r *http.Request) {
	oauthStateString := "test_state"

	// リクエストからstateを取得
	state := r.FormValue("state")
	if state != oauthStateString {
		log.Printf("Invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// 認証コードを取得
	code := r.FormValue("code")
	var err error
	token, err = spotifyConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Code exchange failed with '%v'\n", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// 認証が成功したので、プロフィールページにリダイレクト
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// /profileエンドポイント: ユーザープロフィールと再生履歴を表示
func handleProfile(w http.ResponseWriter, r *http.Request) {
	if token == nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	client := spotifyConfig.Client(context.Background(), token)

	// ユーザープロフィールを取得
	userProfile, err := getUserProfile(client)
	if err != nil {
		log.Printf("Failed to get user profile: %v\n", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// 再生履歴を取得
	recentTracks, err := getRecentlyPlayedTracks(client)
	if err != nil {
		log.Printf("Failed to get recently played tracks: %v\n", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// ユーザープロフィールと再生履歴を表示
	fmt.Fprintf(w, "User Profile: %v\n", userProfile)
	extractTrackInfo(recentTracks, client)
}

// ユーザープロフィールを取得する関数
func getUserProfile(client *http.Client) (map[string]interface{}, error) {
	resp, err := client.Get("https://api.spotify.com/v1/me")
	if err != nil {
		return nil, fmt.Errorf("Failed to get user profile: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Failed to get user profile: %s", body)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read user profile response: %v", err)
	}

	var userProfile map[string]interface{}
	if err := json.Unmarshal(body, &userProfile); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal user profile: %v", err)
	}

	return userProfile, nil
}

// ユーザーの再生履歴を取得する関数
func getRecentlyPlayedTracks(client *http.Client) ([]byte, error) {
	resp, err := client.Get("https://api.spotify.com/v1/me/player/recently-played")
	if err != nil {
		return nil, fmt.Errorf("Failed to get recently played tracks: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Failed to get recently played tracks: %s", body)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read recently played tracks response: %v", err)
	}
	return body, nil
}

// 再生履歴から必要な情報を抽出し、YouTube検索用クエリを作成
func extractTrackInfo(recentTracks []byte, client *http.Client) {

	var result map[string]interface{}
	if err := json.Unmarshal(recentTracks, &result); err != nil {
		log.Fatalf("Failed to unmarshal recently played tracks: %v", err)
	}

	items := result["items"].([]interface{})
	for _, item := range items {
		track := item.(map[string]interface{})["track"].(map[string]interface{})
		trackName := track["name"].(string)   // 曲名
		album := track["album"].(map[string]interface{})
		albumName := album["name"].(string)   // アルバム名
		fmt.Printf("Album Name: %s\n", albumName) // これでアルバム名を表示

		// アルバムの最も小さい画像の取得
		albumImages := album["images"].([]interface{})
		if len(albumImages) > 0 {
			// 一番小さいサイズの画像はリストの最後の要素になることが多い
			smallestImage := albumImages[len(albumImages)-1].(map[string]interface{})["url"].(string)
			fmt.Printf("Album Image (Smallest) URL: %s\n", smallestImage)
		}

		// アーティスト名とその他の情報を1つのリストにまとめる
		artists := track["artists"].([]interface{})
		artistNames := []string{}
		artistSpotifyURLs := []string{}
		artistImageURLs := []string{}

		for _, artist := range artists {
			artistName := artist.(map[string]interface{})["name"].(string)
			artistSpotifyURL := artist.(map[string]interface{})["external_urls"].(map[string]interface{})["spotify"].(string)
			artistNames = append(artistNames, artistName)
			artistSpotifyURLs = append(artistSpotifyURLs, artistSpotifyURL)

			// アーティスト画像の最も小さい画像の取得
			artistID := artist.(map[string]interface{})["id"].(string)
			artistImages := getArtistImages(client, artistID)
			if len(artistImages) > 0 {
				// 一番小さいサイズの画像のみ追加
				artistImageURLs = append(artistImageURLs, artistImages[len(artistImages)-1])
			}
		}

		// アーティスト名の出力
		fmt.Printf("Artists: %s\n", artistNames)

		// アーティストのSpotify URLの出力
		for _, url := range artistSpotifyURLs {
			fmt.Printf("Spotify Artist Page: %s\n", url)
		}

		// アーティスト画像の出力（最も小さい画像のみ）
		for _, imageURL := range artistImageURLs {
			fmt.Printf("Artist Image (Smallest) URL: %s\n", imageURL)
		}

		// YouTube検索クエリの作成（全アーティスト名を連結してクエリに含める）
		searchQuery := fmt.Sprintf("%s %s Official Music Video", trackName, artistNames)
		fmt.Printf("YouTube search query: %s\n", searchQuery)
	}

}

// アーティストの画像を取得する関数
func getArtistImages(client *http.Client, artistID string) []string {
	artistURL := fmt.Sprintf("https://api.spotify.com/v1/artists/%s", artistID)
	resp, err := client.Get(artistURL)
	if err != nil {
		log.Printf("Failed to get artist data: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	artistDataMap := map[string]interface{}{}
	if err := json.Unmarshal(body, &artistDataMap); err != nil {
		log.Printf("Failed to unmarshal artist data: %v\n", err)
		return nil
	}

	artistImages := artistDataMap["images"].([]interface{})
	imageURLs := []string{}
	for _, image := range artistImages {
		imageURL := image.(map[string]interface{})["url"].(string)
		imageURLs = append(imageURLs, imageURL)
	}
	return imageURLs
}
