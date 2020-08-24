package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Execute an HTTP request with a response
func ExecuteRequest(request *http.Request) (map[string]interface{}, error) {
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
		panic(err)
	}

	// Convert JSON to map
	response := responseContainer.(map[string]interface{})
	return response, nil
}

// Generic HTTP client
func HttpClient() *http.Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	return client
}
