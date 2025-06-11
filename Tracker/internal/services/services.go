package services

import (
	"Tracker/internal/config"
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// AIService handles AI-related operations
type AIService struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewAIService creates a new AI service instance
func NewAIService() (*AIService, error) {
	ctx := context.Background()

	// Get Gemini API key from config
	apiKey := config.GetGeminiApiKey()
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini API key not found")
	}

	// Create Gemini client
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %v", err)
	}

	// Get model from config or use default
	modelName := config.GetGeminiApiKey()
	if modelName == "" {
		modelName = "gemini-pro"
	}

	// Initialize the model
	model := client.GenerativeModel(modelName)

	return &AIService{
		client: client,
		model:  model,
	}, nil
}

// GetActivitySuggestions generates activity suggestions based on user preferences
func (s *AIService) GetActivitySuggestions(preferences string) ([]string, error) {
	ctx := context.Background()

	prompt := "Based on these preferences: " + preferences + "\nSuggest 5 activities that would be suitable. Format each suggestion as a single line."

	// Generate content
	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %v", err)
	}

	// Get the response text
	response := resp.Candidates[0].Content.Parts[0].(genai.Text)

	// Split the response into individual suggestions
	suggestions := strings.Split(string(response), "\n")

	// Clean up suggestions (remove empty lines and numbers)
	cleanSuggestions := make([]string, 0)
	for _, suggestion := range suggestions {
		suggestion = strings.TrimSpace(suggestion)
		if suggestion != "" {
			// Remove leading numbers and dots if present
			suggestion = strings.TrimLeft(suggestion, "0123456789. ")
			cleanSuggestions = append(cleanSuggestions, suggestion)
		}
	}

	return cleanSuggestions, nil
}

// Close closes the AI service client
func (s *AIService) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}
