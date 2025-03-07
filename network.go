package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ModelsData contains the models array
type ModelsData struct {
	Models []string `json:"models"`
}

// ListModelsResponse represents the response from the models endpoint
type ListModelsResponse struct {
	Data    ModelsData `json:"data"`
	Message string     `json:"message"`
	Status  int        `json:"status"`
}

// Message represents a single message in a conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionRequest represents the request payload for the chat completions API
type CompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// ResponseChunk represents a single chunk of the streaming response
type ResponseChunk struct {
	Model     string  `json:"model"`
	CreatedAt string  `json:"created_at"`
	Message   Message `json:"message"`
	Done      bool    `json:"done"`
}

// createCompletionRequest creates a chat completion request payload
func createCompletionRequest(model string, messages []Message) ([]byte, error) {
	requestBody := CompletionRequest{
		Model:    model,
		Messages: messages,
	}

	return json.Marshal(requestBody)
}

// sendRequest sends an HTTP request to the API
func sendRequest(url string, jsonData []byte, apiKey string) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	return client.Do(req)
}

// processStreamResponse processes the streaming response and concatenates content
func processStreamResponse(body io.Reader) (string, error) {
	var fullContent strings.Builder
	scanner := bufio.NewScanner(body)
	jobID := ""

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Handle event lines
		if strings.HasPrefix(line, "event:") {
			continue
		}

		// Handle data lines
		if strings.HasPrefix(line, "data:") {
			dataContent := strings.TrimSpace(strings.TrimPrefix(line, "data:"))

			// Try to parse as JSON
			var chunk ResponseChunk
			if err := json.Unmarshal([]byte(dataContent), &chunk); err != nil {
				if jobID == "" {
					jobID = dataContent // Assume first non-JSON data is job ID
				}
			} else {
				// Successfully parsed - add to content
				fullContent.WriteString(chunk.Message.Content)
			}
			continue
		}

		// Try to parse direct JSON (no data: prefix)
		var chunk ResponseChunk
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			fmt.Printf("Warning: Couldn't parse line as JSON: %s\n", line)
			continue
		} else {
			fullContent.WriteString(chunk.Message.Content)
		}
	}

	if err := scanner.Err(); err != nil {
		return fullContent.String(), fmt.Errorf("error reading response: %v", err)
	}

	if jobID != "" {
		fmt.Printf("[Job ID: %s]\n", jobID)
	}

	return fullContent.String(), nil
}

// ListAvailableModels fetches the list of available models from the API
func ListAvailableModels(apiKey string) ([]string, error) {
	url := "https://anura-testnet.lilypad.tech/api/v1/models"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("received status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var modelsResp ListModelsResponse
	if err := json.Unmarshal(body, &modelsResp); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	// Return the list of models
	return modelsResp.Data.Models, nil
}
