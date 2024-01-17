package main

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type CompletionServiceInterface interface {
	GetCompletion(diff, model, apiKey, prompt string) (string, error)
}

type OpenAICompletionService struct{}

func (s *OpenAICompletionService) GetCompletion(diff, gptModel, apiKey, prompt string) (string, error) {

	fullPrompt := fmt.Sprintf(prompt, diff)

	// OpenAI
	client := openai.NewClient(apiKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: gptModel,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fullPrompt,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
