package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	//"os"
	"io/ioutil"
	"encoding/json"

	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/spotify"

	//"github.com/joho/godotenvs"
)

func HandleUserProfile(w http.ResponseWriter, r *http.Request){
	if token == nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	client := spotifyConfig.Client(context.Background(), token)
	userProfile, err := GetUserProfile(client)
	if err != nil {
		log.Printf("Failed to get recently played tracks: %v\n", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	recentlyPlayedTracks, err := GetUserRecentPlayed(client)
	if err != nil {
		log.Printf("Failed to get recently played tracks: %v\n", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	sendData := map[string]interface{}{
		"user_profile": userProfile,
		"recently_played_tracks": recentlyPlayedTracks,
	}

    w.Header().Set("Content-Type", "application/json")

    // マップをJSONにエンコードしてレスポンスとして送信
    json.NewEncoder(w).Encode(sendData)
}



//userProfileとrecentPlaylistをまとめる
// func connectData(client *http.Client)(map[string] interface{}){
// 	profileBody, err := GetUserProfile(client)
// 	if err != nil {
// 		fmt.Println("Error fetching user profile:", err)
// 		return nil 
// 	}
// 	userProfile := ExtrackUserProfile(client, profileBody)
// 	recentlyBody, err := GetUserRecentPlayed(client)
// 	if err != nil {
// 		fmt.Println("Error fetching recently played tracks:", err)
// 		return nil 
// 	}
// 	recentlyPlayedTracks := ExtrackFromRecentPlayed(client, recentlyBody)
// 	sendData := map[string]interface{}{
// 		"user_profile": userProfile,
// 		"recently_played_tracks": recentlyPlayedTracks,
// 	}
// 	return sendData
// }



//Spotify履歴の取得

func GetUserProfile(client *http.Client)(map[string]string, error){
	resp, err := client.Get("https://api.spotify.com/v1/me")

	if err != nil{
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}

	userProfile, err := ExtrackUserProfile(client, body)
	if err != nil{
		return nil, err
	}
	return userProfile, nil
}

func ExtrackUserProfile(client *http.Client, body []byte)(map[string]string, error){
	userProfile := map[string]string{}
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal user profile: %v", err)
	}

	userProfile["id"] = result["id"].(string)
	userProfile["display_name"] = result["display_name"].(string)
	images := result["images"].([]interface{})
	if len(images) > 0{
		firstImage := images[0].(map[string]interface{})
		userProfile["url"] = firstImage["url"].(string)
	} else {
		userProfile["url"] = ""
	}
	userProfile["country"] = result["country"].(string)
	return userProfile, nil
}

func GetUserRecentPlayed(client *http.Client)([]interface{}, error){
	resp, err := client.Get("https://api.spotify.com/v1/me/player/recently-played")

	if err != nil{
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}

	recentlyPlayedTracks, err := ExtrackFromRecentPlayed(client, body)
	if err != nil{
		return nil, err
	}
	return recentlyPlayedTracks, nil
}

//履歴の整形
func ExtrackFromRecentPlayed(client *http.Client, body []byte)([]interface{}, error){
	var result map[string] interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		// log.Printf("Failed to unmarshal user profile: %s\n", err)
		// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil, fmt.Errorf("Failed to unmarshal user profile: %v", err)
	}
	items := result["items"].([]interface{})
	recentlyPlayedTracks := []interface{}{}
	for _, item := range items{
		track := item.(map[string]interface{})["track"].(map[string]interface{})
		album := track["album"].(map[string]interface{})
		trackName := track["name"]

		images := album["images"].([]interface{})
		albumImageUrl := ""
		if len(images) > 0{
			firstAlbumImageUrl := images[0].(map[string]interface{})
			albumImageUrl = firstAlbumImageUrl["url"].(string)			
		}



		albumName := album["name"].(string)
		artists := album["artists"].([]interface{})
		artistsData := []interface{}{}
		for _, artist := range artists{
			artistMap := artist.(map[string]interface{})
			urls := artistMap["external_urls"].(map[string]interface{})
			spotifyUrl := urls["spotify"].(string)
			artistName := artistMap["name"].(string)
			artistId := artistMap["id"].(string)
			artistImageUrl := GetArtistsImageUrl(client, artistId)
			artistData := map[string]interface{}{
				"name": artistName,
				"spotify_url": spotifyUrl,
				"smallest_image_url": artistImageUrl,
			}
			artistsData = append(artistsData, artistData)
		}
		trackData := map[string]interface{}{
			"track_name": trackName,
			"album": map[string]string{
				"name": albumName,
				"smallest_image_url": albumImageUrl,
			},
			"artists": artistsData,
			//Youtube検索クエリ
		}
		recentlyPlayedTracks = append(recentlyPlayedTracks, trackData)
	}
	return recentlyPlayedTracks, nil
}

//idからartistsのimgaeUrlをとってくる
func GetArtistsImageUrl(client *http.Client, id string)(string){
	url := fmt.Sprintf("https://api.spotify.com/v1/artists/%s", id)
	resp, err := client.Get(url)
	if err != nil{
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return ""
	}

	var result map[string] interface{}
	if err := json.Unmarshal(body, &result); err != nil{
		return ""
	}
	images := result["images"].([]interface{})
	imageUrl := ""
	if len(images) > 0{
		firstImage := images[0].(map[string]interface{})
		imageUrl = firstImage["url"].(string)		
	}

	return imageUrl
}
