package api

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// TwitchRequestData Twitch API & account data
type TwitchRequestData struct {
	EndpointAuth   string
	EndpointStream string
	ClientId       string
	ClientSecret   string
	UserId         string
	UserName       string
}

// Populate API data to make Twitch
func getTwitchRequestData() TwitchRequestData {
	twitchRequestData := TwitchRequestData{
		EndpointAuth:   "https://id.twitch.tv/oauth2/token",
		EndpointStream: "https://api.twitch.tv/helix/streams",
		ClientId:       os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret:   os.Getenv("TWITCH_CLIENT_SECRET"),
		UserId:         os.Getenv("TWITCH_USER_ID"),
		UserName:       os.Getenv("TWITCH_USER_NAME"),
	}
	return twitchRequestData
}

// Get Twitch bearer token
func getTwitchToken() string {
	tokenReq := getTokenRequest()
	tokenRes := ExecuteRequest(tokenReq)
	token := tokenRes["access_token"].(string)
	return token
}

// Construct HTTP request to fetch Twitch token
func getTokenRequest() *http.Request {
	twitchData := getTwitchRequestData()
	requestUrl := fmt.Sprintf(
		"%s?grant_type=client_credentials&client_id=%s&client_secret=%s",
		twitchData.EndpointAuth,
		twitchData.ClientId,
		twitchData.ClientSecret)
	req, reqErr := http.NewRequest("POST", requestUrl, nil)
	if reqErr != nil {
		log.Fatalf("Unable to get Twitch token: %v", reqErr.Error())
	}
	return req
}

// Fetch live stream(s) by user login
func getLiveStreamByUserRequest() *http.Request {
	twitchData := getTwitchRequestData()
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

// GetTwitchStream Check if Twitch stream is live
func GetTwitchStream() bool {
	streamReq := getLiveStreamByUserRequest()
	streamRes := ExecuteRequest(streamReq)
	streamData := streamRes["data"]
	if streamData != nil {
		log.Print("Twitch steam is live!")
		return true
	}
	log.Print("Twitch steam currently offline.")
	return false
}
