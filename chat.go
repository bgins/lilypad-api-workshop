package main

import (
	"fmt"
	"io"
	"net/http"
)

// CompleteChat makes a request to the AI service and returns the generated response
func CompleteChat(apiKey, model string, messages []Message) (string, error) {
	// Create the request payload
	jsonData, err := createCompletionRequest(model, messages)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Send the request
	url := "https://anura-testnet.lilypad.tech/api/v1/chat/completions"
	resp, err := sendRequest(url, jsonData, apiKey)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("received status code %d: %s", resp.StatusCode, string(body))
	}

	// Process the response
	return processStreamResponse(resp.Body)
}
