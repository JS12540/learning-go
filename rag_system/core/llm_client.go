package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"rag_system/models"
)

var httpClient = &http.Client{}

func GenerateChatCompletion(messages []models.ChatCompletionMessage, modelName string) (string, error) {
	if modelName == "" {
		modelName = "gpt-4.1-mini"
	}

	reqPayload := models.ChatCompletionRequest{
		Model:    modelName,
		Messages: messages,
	}

	payloadBytes, err := json.Marshal(reqPayload)

	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	apiURL := "https://api.openai.com/v1/chat/completions"

	req, err := http.NewRequest(
		"POST",
		apiURL,
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API error: %s", string(bodyBytes))
	}

	var completionResponse models.ChatCompletionResponse

	err = json.Unmarshal(bodyBytes, &completionResponse)
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(completionResponse.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return completionResponse.Choices[0].Message.Content, nil
}
