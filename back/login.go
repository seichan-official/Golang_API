package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	//"os"
	// "io/ioutil"
	

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"

	//"github.com/joho/godotenv"
)

err := godotenv.Load("./.env")
if err != nil{
	log.Fatalf("Error load env faile: %v", err)
}

// Spotify APIのOAuth2設定
var spotifyConfig = &oauth2.Config{



	// ClientID:os.Getenv("client_id"),// Spotify Developerから取得
	// ClientSecret: os.Getenv("seacret_id"),    // Spotify Developerから取得


	RedirectURL:  "http://localhost:8080/", // リダイレクトURL
	Endpoint:     spotify.Endpoint,        // Spotify用のOAuth2エンドポイント
	Scopes: []string{
		"user-read-private", "user-read-email",
		"user-read-recently-played", // 必要なスコープを指定
	},
}

var oauthStateString = "random"
var token *oauth2.Token



// メイン関数でサーバーを開始
func main() {


	http.HandleFunc("/api/spotify/login", handleLogin)
	http.HandleFunc("/", handleCallback)
	http.HandleFunc("/api/spotify/history", HandleUserProfile)

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
	localToken, err := spotifyConfig.Exchange(context.Background(), code)
	
	if err != nil {
		log.Printf("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	token = localToken
	//client := spotifyConfig.Client(context.Background(), token)

	http.Redirect(w, r, "/api/spotify/history", http.StatusSeeOther)

}
