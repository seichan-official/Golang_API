package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	//"os"
	"io/ioutil"
	"encoding/json"
	"Golang_API/model"
	

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"

	//"github.com/joho/godotenvs"
)

// Spotify APIのOAuth2設定
var spotifyConfig = &oauth2.Config{



	// ClientID:os.Getenv("client_id"),// Spotify Developerから取得
	// ClientSecret: os.Getenv("seacret_id"),    // Spotify Developerから取得


	RedirectURL:  "http://localhost:8080/", // リダイレクトURL
	Endpoint:     spotify.Endpoint,        // Spotify用のOAuth2エンドポイント
	Scopes: []string{
		"user-read-private", "user-read-email", // 必要なスコープを指定
	},
}

var oauthStateString = "random"
var userID = ""


// メイン関数でサーバーを開始
func main() {

	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/", handleCallback)

	fmt.Println("Server started at http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// /loginエンドポイント: Spotifyの認証ページにリダイレクト
func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := spotifyConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// /callbackエンドポイント: 認証後に呼び出される
func handleCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {

		log.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := spotifyConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	client := spotifyConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.spotify.com/v1/me")
	if err != nil {
		log.Printf("Failed to get user profile: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	log.Println("Response body: ", string(body))

	// ユーザープロフィール情報をJSONとして解析して表示
    var userProfile struct {
        ID          string `json:"id"`
        DisplayName string `json:"display_name"`
        Email       string `json:"email"`
    }
	

	playlists, err := getUserPlaylists(client, userID)
	if err := json.Unmarshal(body, &userProfile); err != nil {
		log.Printf("Failed to unmarshal user profile: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	var playlistsMap map[string] interface{}
	if err := json.Unmarshal(playlists, &playlistsMap); err != nil {
		log.Printf("Failed to unmarshal user profile: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//dbInsert(userProfile.ID, userProfle.DisplayName, userProfile.Email)
	fmt.Fprintf(w, "User Profile: %+v\n", userProfile)
}

//プレイリストの取得
func getUserPlaylists(client *http.Client, userID string)([]byte, error){
	url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", userID)
	resp, err := client.Get(url)
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}

	return body, nil
}