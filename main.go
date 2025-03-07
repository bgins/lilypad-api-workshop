package main

import (
	"fmt"
	"os"
)

func main() {
	// Get API key
	apiKey, err := getAPIKey()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Get model from environment variable
	model := getModel()
	fmt.Printf("Using model: %s\n", model)

	// Prepare messages
	messages := []Message{
		{Role: "system", Content: "You are a helpful assistant."},
		{Role: "assistant", Content: "Hi how can I help you today"},
		{Role: "user", Content: "What is AI engineering?"},
	}

	// Make the API call
	content, err := CompleteChat(apiKey, model, messages)
	if err != nil {
		fmt.Printf("Error completing chat: %v\n", err)
		return
	}

	// Print the complete message
	fmt.Println("\nComplete response:")
	fmt.Println(content)
}

// getAPIKey gets the API key from environment variables
func getAPIKey() (string, error) {
	apiKey := os.Getenv("ANURA_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("ANURA_API_KEY environment variable is not set")
	}
	return apiKey, nil
}

// getModel gets the model name from environment variables or returns a default
func getModel() string {
	model := os.Getenv("ANURA_MODEL")
	if model == "" {
		// Default model if not specified
		return "qwen2.5:7b"
	}
	return model
}
