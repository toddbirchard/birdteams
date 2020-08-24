package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

// Twitch API & account data
type TwitchData struct {
	EndpointAuth   string
	EndpointStream string
	ClientId       string
	ClientSecret   string
	UserId         string
}

// Load environment variables
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// Populate Twitch user's API data
func GetTwitchData() TwitchData {
	LoadEnv()
	twitchData := TwitchData{
		EndpointAuth:   "https://id.twitch.tv/oauth2/token",
		EndpointStream: "https://api.twitch.tv/helix/streams",
		ClientId:       os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret:   os.Getenv("TWITCH_CLIENT_SECRET"),
		UserId:         os.Getenv("TWITCH_USER_ID"),
	}
	return twitchData
}

// Get Twitch bearer token
func GetTwitchToken() string {
	tokenReq, _ := GetTokenRequest()
	tokenRes, _ := ExecuteRequest(tokenReq)
	token := tokenRes["access_token"].(string)
	log.Println(token)
	return token
}

// Construct HTTP request to fetch Twitch token
func GetTokenRequest() (*http.Request, error) {
	twitchData := GetTwitchData()
	requestUrl := fmt.Sprintf(
		"%s?grant_type=client_credentials&client_id=%s&client_secret=%s",
		twitchData.EndpointAuth,
		twitchData.ClientId,
		twitchData.ClientSecret)
	req, reqErr := http.NewRequest("POST", requestUrl, nil)
	if reqErr != nil {
		return nil, reqErr
	}
	return req, nil
}
