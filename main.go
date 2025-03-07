package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

func main() {
	fmt.Println("Welcome to the Anura Chat CLI!")
	fmt.Println("─────────────────────────────────────────────────────────────")

	// Get API key
	apiKey, err := getAPIKey()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Fetch available models for selection menu
	fmt.Println("Fetching available models...")
	models, err := ListAvailableModels(apiKey)
	if err != nil {
		fmt.Printf("Error fetching models: %v\n", err)
		// Use a default model if fetching fails
		fmt.Println("Using default model: qwen2.5:7b")
		models = []string{"qwen2.5:7b"}
	}

	// Create a prompt select using promptui
	prompt := promptui.Select{
		Label: "Select AI Model",
		Items: models,
		Size:  10, // Show 10 items at once
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "➤ {{ . | green }}",
			Inactive: "  {{ . }}",
			Selected: "✓ Selected model: {{ . | green }}",
		},
	}

	_, model, err := prompt.Run()
	if err != nil {
		fmt.Printf("Selection failed: %v\n", err)
		return
	}

	// Initialize conversation with system message
	messages := []Message{
		{Role: "system", Content: "You are a helpful assistant."},
	}

	// Create a scanner for reading user input
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("\nUsing model: %s\n", model)
	fmt.Println("Type your messages and press Enter.")
	fmt.Println("Type 'exit', 'quit', or press Ctrl+C to end the conversation.")
	fmt.Println("─────────────────────────────────────────────────────────────")

	// Main conversation loop
	for {
		// Prompt for user input
		fmt.Print("\nYou: ")
		if !scanner.Scan() {
			break
		}

		userInput := strings.TrimSpace(scanner.Text())
		if userInput == "" {
			continue
		}

		// Check for exit commands
		if userInput == "exit" || userInput == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		// Add user message to the conversation
		messages = append(messages, Message{Role: "user", Content: userInput})

		// Call Anura
		fmt.Println("\nAnura is processing...")
		response, err := CompleteChat(apiKey, model, messages)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Display Anura response
		fmt.Printf("\nAI: %s\n", response)

		// Add Anura response to the conversation history
		messages = append(messages, Message{Role: "assistant", Content: response})
	}
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
