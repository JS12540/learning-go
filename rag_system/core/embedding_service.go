package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"rag_system/config"
	"rag_system/models"
	"strings"
	"time"
)

var NewhttpClient = &http.Client{Timeout: 180 * time.Second}

const (
	defaultEmbeddingBatchSize = 32
	maxTokensPerBatch         = 8000
	maxCharsPerToken          = 4
	maxBatchSizeLimit         = 64
	minBatchSize              = 1
)

const openAIEmbeddingURL = "https://api.openai.com/v1/embeddings"

func GetEmbeddings(texts []string, modelName string) ([][]float32, error) {

	if modelName == "" {
		modelName = config.AppConfig.EmbeddingModel
	}

	if len(texts) == 0 {
		return [][]float32{}, nil
	}

	allEmbeddings := make([][]float32, len(texts))

	batches := createAdaptiveBatches(texts)

	log.Printf("Processing %d texts in %d adaptive batches", len(texts), len(batches))

	for batchIndex, batch := range batches {

		embeddings, err := processBatchWithRetry(batch, modelName, batchIndex)
		if err != nil {
			return nil, fmt.Errorf("failed to process batch %d: %w", batchIndex, err)
		}

		for i, embedding := range embeddings {

			globalIndex := batch.StartIndex + i

			if globalIndex < len(allEmbeddings) {
				allEmbeddings[globalIndex] = embedding
			}

		}

		log.Printf("Successfully processed batch %d (%d texts)", batchIndex, len(batch.Texts))
	}

	for idx, emb := range allEmbeddings {

		if len(emb) == 0 {
			return nil, fmt.Errorf("embedding for text at index %d was not populated", idx)
		}

	}

	return allEmbeddings, nil
}

type EmbeddingBatch struct {
	Texts      []string
	StartIndex int
	TotalChars int
}

func createAdaptiveBatches(texts []string) []EmbeddingBatch {

	var batches []EmbeddingBatch

	i := 0

	for i < len(texts) {

		batch := EmbeddingBatch{
			StartIndex: i,
		}

		currentChars := 0
		batchSize := 0

		for i+batchSize < len(texts) && batchSize < maxBatchSizeLimit {

			textChars := len(texts[i+batchSize])

			estimatedTokens := (currentChars + textChars) / maxCharsPerToken

			if estimatedTokens > maxTokensPerBatch && batchSize > 0 {
				break
			}

			if textChars/maxCharsPerToken > maxTokensPerBatch {

				log.Printf(
					"Warning: Text at index %d is very large (%d chars, ~%d tokens), processing individually",
					i+batchSize,
					textChars,
					textChars/maxCharsPerToken,
				)

				if batchSize == 0 {
					batch.Texts = append(batch.Texts, texts[i+batchSize])
					batch.TotalChars = textChars
					batchSize = 1
				}

				break
			}

			batch.Texts = append(batch.Texts, texts[i+batchSize])
			currentChars += textChars
			batchSize++

		}

		batch.TotalChars = currentChars
		batches = append(batches, batch)

		i += batchSize
	}

	return batches
}

func getEmbeddingDimension(modelName string) int {

	modelDimensions := map[string]int{
		"text-embedding-3-small": 1536,
		"text-embedding-3-large": 3072,
		"text-embedding-ada-002": 1536,
	}

	if dim, exists := modelDimensions[modelName]; exists {
		return dim
	}

	log.Printf("Unknown model %s, defaulting to 1536", modelName)

	return 1536
}

func processBatchWithRetry(batch EmbeddingBatch, modelName string, batchIndex int) ([][]float32, error) {

	currentBatch := batch

	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {

		log.Printf(
			"Batch %d attempt %d: %d texts, %d chars (~%d tokens)",
			batchIndex,
			attempt+1,
			len(currentBatch.Texts),
			currentBatch.TotalChars,
			currentBatch.TotalChars/maxCharsPerToken,
		)

		embeddings, err := sendEmbeddingRequest(currentBatch.Texts, modelName)

		if err == nil {
			return embeddings, nil
		}

		if isOversizedBatchError(err) {

			if len(currentBatch.Texts) == 1 {

				log.Printf(
					"Single text at batch %d is too large (%d chars), skipping",
					batchIndex,
					currentBatch.TotalChars,
				)

				dimension := getEmbeddingDimension(modelName)

				placeholder := make([]float32, dimension)

				return [][]float32{placeholder}, nil
			}

			if len(currentBatch.Texts) > minBatchSize {

				log.Printf(
					"Batch %d is too large, splitting in half (attempt %d)",
					batchIndex,
					attempt+1,
				)

				midpoint := len(currentBatch.Texts) / 2

				firstHalfChars := 0

				for _, text := range currentBatch.Texts[:midpoint] {
					firstHalfChars += len(text)
				}

				secondHalfChars := 0

				for _, text := range currentBatch.Texts[midpoint:] {
					secondHalfChars += len(text)
				}

				firstHalf := EmbeddingBatch{
					Texts:      currentBatch.Texts[:midpoint],
					StartIndex: currentBatch.StartIndex,
					TotalChars: firstHalfChars,
				}

				secondHalf := EmbeddingBatch{
					Texts:      currentBatch.Texts[midpoint:],
					StartIndex: currentBatch.StartIndex + midpoint,
					TotalChars: secondHalfChars,
				}

				firstEmbeddings, err1 := processBatchWithRetry(firstHalf, modelName, batchIndex)

				if err1 != nil {
					return nil, fmt.Errorf("failed to process first half of split batch: %w", err1)
				}

				secondEmbeddings, err2 := processBatchWithRetry(secondHalf, modelName, batchIndex)

				if err2 != nil {
					return nil, fmt.Errorf("failed to process second half of split batch: %w", err2)
				}

				combined := append(firstEmbeddings, secondEmbeddings...)

				return combined, nil
			}
		}

		if attempt == maxRetries-1 || len(currentBatch.Texts) <= minBatchSize {

			return nil, fmt.Errorf("failed after %d attempts: %w", attempt+1, err)
		}

		time.Sleep(time.Second * time.Duration(attempt+1))
	}

	return nil, fmt.Errorf("exceeded maximum retry attempts")
}

func sendEmbeddingRequest(texts []string, modelName string) ([][]float32, error) {

	reqPayload := models.EmbeddingRequest{
		Input: texts,
		Model: modelName,
	}

	payloadBytes, err := json.Marshal(reqPayload)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal embedding request: %w", err)
	}

	req, err := http.NewRequest("POST", openAIEmbeddingURL, bytes.NewBuffer(payloadBytes))

	if err != nil {
		return nil, fmt.Errorf("failed to create embedding request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set(
		"Authorization",
		"Bearer "+os.Getenv("OPENAI_API_KEY"),
	)

	resp, err := NewhttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI embedding API: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		errBodyBytes, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"embedding API request failed with status %s: %s",
			resp.Status,
			string(errBodyBytes),
		)
	}

	var embeddingResp models.EmbeddingAPIResponse

	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
		return nil, fmt.Errorf("failed to decode embedding API response: %w", err)
	}

	if len(embeddingResp.Data) != len(texts) {

		return nil, fmt.Errorf(
			"mismatch in number of embeddings returned (%d) vs texts sent (%d)",
			len(embeddingResp.Data),
			len(texts),
		)
	}

	embeddings := make([][]float32, len(texts))

	for _, data := range embeddingResp.Data {

		if data.Index >= 0 && data.Index < len(embeddings) {

			embeddings[data.Index] = data.Embedding

		} else {

			return nil, fmt.Errorf(
				"embedding data index out of bounds: %d",
				data.Index,
			)
		}

	}

	return embeddings, nil
}

func isOversizedBatchError(err error) bool {

	errorStr := strings.ToLower(err.Error())

	oversizedIndicators := []string{
		"too large",
		"input is too large",
		"context length exceeded",
		"maximum context length",
		"token limit",
		"input size",
	}

	for _, indicator := range oversizedIndicators {

		if strings.Contains(errorStr, indicator) {
			return true
		}

	}

	return false
}
