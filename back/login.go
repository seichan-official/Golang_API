package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

// Spotify APIのOAuth2設定
var spotifyConfig = &oauth2.Config{
	
	ClientID:"",// Spotify Developerから取得
	ClientSecret: "",    // Spotify Developerから取得

	RedirectURL:  "http://localhost:8080/callback", // リダイレクトURL
	Endpoint:     spotify.Endpoint,        // Spotify用のOAuth2エンドポイント
	Scopes: []string{
		"user-read-private", "user-read-email", // 必要なスコープを指定
	},
}

var oauthStateString = "random"

// メイン関数でサーバーを開始
func main() {
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)

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

	fmt.Fprintf(w, "Access Token: %s\n", token.AccessToken)
}