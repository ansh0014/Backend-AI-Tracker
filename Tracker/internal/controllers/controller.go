package services

import (
	"Tracker/internal/config"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// AIService handles AI-related operations
type AIService struct {
	client *genai.Client
	model  *genai.GenerativeModel
	ctx    context.Context
}

// NewAIService creates a new AI service instance with timeout context
func NewAIService(timeout time.Duration) (*AIService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	apiKey := config.GetGeminiApiKey()
	if apiKey == "" {
		return nil, fmt.Errorf("gemini API key not found")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %v", err)
	}

	modelName := config.GetGeminiModel() // Add this to config package
	if modelName == "" {
		modelName = "gemini-pro"
	}

	model := client.GenerativeModel(modelName)
	model.SetTemperature(0.7) // Add some creativity while keeping responses focused

	return &AIService{
		client: client,
		model:  model,
		ctx:    context.Background(),
	}, nil
}

// GetActivitySuggestions generates activity suggestions with context and safety checks
func (s *AIService) GetActivitySuggestions(ctx context.Context, preferences string) ([]string, error) {
	if preferences == "" {
		return nil, fmt.Errorf("preferences cannot be empty")
	}

	prompt := fmt.Sprintf(`Based on these preferences: %s
Generate 5 specific activities that would be suitable.
Format each suggestion as a clear, actionable item on a new line.
Focus on productive and healthy activities.`, preferences)

	resp, err := s.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("content generation failed: %v", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return nil, fmt.Errorf("no valid response generated")
	}

	response := resp.Candidates[0].Content.Parts[0].(genai.Text)
	suggestions := strings.Split(string(response), "\n")

	return cleanSuggestions(suggestions), nil
}

// cleanSuggestions helper function to process suggestions
func cleanSuggestions(raw []string) []string {
	clean := make([]string, 0, len(raw))
	for _, suggestion := range raw {
		suggestion = strings.TrimSpace(suggestion)
		if suggestion == "" {
			continue
		}
		// Remove leading numbers, dots, and dashes
		suggestion = strings.TrimLeft(suggestion, "0123456789.-• ")
		clean = append(clean, suggestion)
	}
	return clean
}

// Close safely closes the AI service client with context
func (s *AIService) Close(ctx context.Context) error {
	if s.client != nil {
		errChan := make(chan error, 1)
		go func() {
			errChan <- s.client.Close()
		}()

		select {
		case err := <-errChan:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return nil
}