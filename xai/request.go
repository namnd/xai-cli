package xai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

const (
	API_ENDPOINT = "https://api.x.ai/v1/chat/completions"
)

func MakeAPICall(ctx context.Context, apiKey string, payload []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", API_ENDPOINT, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return body, nil
}
