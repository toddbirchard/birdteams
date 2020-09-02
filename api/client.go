package api

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Load environment variables
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// Execute an HTTP request with a response
func ExecuteRequest(request *http.Request) map[string]interface{} {
	// Construct HTTP request to fetch session data.
	client := HttpClient()
	res, resError := client.Do(request)
	if resError != nil {
		log.Fatal(resError)
	}

	// Parse response data
	data, dataErr := ioutil.ReadAll(res.Body)
	if dataErr != nil {
		panic(dataErr)
	}

	// Render response as JSON
	var responseContainer interface{}
	err := json.Unmarshal(data, &responseContainer)
	if err != nil {
		log.Fatalf("Unable to parse JSON: %v", err)
		return nil
	}

	// Convert JSON to map
	response := responseContainer.(map[string]interface{})
	return response
}

// Generic HTTP client
func HttpClient() *http.Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	return client
}
