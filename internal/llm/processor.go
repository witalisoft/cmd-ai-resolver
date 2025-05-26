package llm

import (
	"cmd-ai-resolver/internal/logger"
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

const (
	defaultOpenAIModel = "gpt-4.1-mini"
)

func ProcessWithOpenAI(fileContent string, extractedPrompt string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		logger.Log.Errorf("OPENAI_API_KEY environment variable not set")
		return "", fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	baseURL := os.Getenv("OPENAI_BASE_URL") // Optional

	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
		logger.Log.Debugf("Using custom OpenAI base URL: %s", baseURL)
	} else {
		logger.Log.Debugf("Using default OpenAI base URL.")
	}

	client := openai.NewClientWithConfig(config)

	model := os.Getenv("OPENAI_BASE_MODEL")
	if model == "" {
		model = defaultOpenAIModel
	}
	logger.Log.Debugf("Using OpenAI model: %s", model)

	llmPrompt := fmt.Sprintf(
		"Given the following shell command:\n\n```shell\n%s\n```\n\nTranslate the AI instruction `\"%s\"` into a shell command segment that achieves the described task. "+
			"The AI instruction is embedded in the command. "+
			"Return *only* the resulting shell command segment, not the entire modified command. "+
			"For example, if the input command is 'ls -l | <AI>show only last 5 lines</AI>' and the instruction is 'show only last 5 lines', "+
			"you should return 'tail -n 5'.",
		fileContent, extractedPrompt,
	)

	logger.Log.Debugf("Sending prompt to LLM: %s", llmPrompt)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: llmPrompt,
				},
			},
		},
	)

	if err != nil {
		logger.Log.Errorf("OpenAI API request failed: %v", err)
		return "", fmt.Errorf("OpenAI API request failed: %w", err)
	}

	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		logger.Log.Errorf("Received an empty response or no choices from LLM")
		return "", fmt.Errorf("received an empty response or no choices from LLM")
	}

	processedSegment := resp.Choices[0].Message.Content
	logger.Log.Debugf("LLM response (processed segment): %s", processedSegment)
	return processedSegment, nil
}
