package api

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
	UserName       string
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
		UserName:       os.Getenv("TWITCH_USER_NAME"),
	}
	return twitchData
}

// Get Twitch bearer token
func GetTwitchToken() string {
	tokenReq := GetTokenRequest()
	tokenRes := ExecuteRequest(tokenReq)
	token := tokenRes["access_token"].(string)
	log.Println(token)
	return token
}

// Fetch live stream(s) by user login
func GetStreamByUser() *http.Request {
	twitchData := GetTwitchData()
	requestUrl := fmt.Sprintf(
		"%s?user_login=%s",
		twitchData.EndpointStream,
		twitchData.UserName)
	req, reqErr := http.NewRequest("GET", requestUrl, nil)
	if reqErr != nil {
		log.Fatalf("Unable to get Twitch stream by user: %v", reqErr)
	}
	return req
}

// Construct HTTP request to fetch Twitch token
func GetTokenRequest() *http.Request {
	twitchData := GetTwitchData()
	requestUrl := fmt.Sprintf(
		"%s?grant_type=client_credentials&client_id=%s&client_secret=%s",
		twitchData.EndpointAuth,
		twitchData.ClientId,
		twitchData.ClientSecret)
	req, reqErr := http.NewRequest("POST", requestUrl, nil)
	if reqErr != nil {
		log.Fatalf("Unable to get Twitch token: %v", reqErr)
	}
	return req
}
