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
	ctx    context.Context
}

// NewAIService creates a new AI service instance
func NewAIService() (*AIService, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	ctx := context.Background()
	client, modelName, err := InitializeGeminiClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Gemini client: %v", err)
	}

	model := client.GenerativeModel(modelName)

	return &AIService{
		client: client,
		model:  model,
		ctx:    ctx,
	}, nil
}


// InitializeGeminiClient sets up the Gemini AI client using config.
func InitializeGeminiClient(ctx context.Context, cfg *config.Config) (*genai.Client, string, error) {
	if cfg.GeminiApiKey == "" {
		return nil, "", fmt.Errorf("Gemini API key not found in configuration")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiApiKey))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create Gemini client: %v", err)
	}

	model := cfg.GeminiModel
	if model == "" {
		model = "gemini-pro"
	}

	return client, model, nil
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
