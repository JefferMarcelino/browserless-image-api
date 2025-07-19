package fetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func FetchContent(ctx context.Context, pageURL string) (string, error) {
	payload := map[string]any{
		"url":         pageURL,
		"gotoOptions": map[string]any{"waitUntil": "networkidle0"},
		"waitForSelector": map[string]any{
			"selector": "body",
			"timeout":  10000,
		},
		"bestAttempt": true,
	}
	body, _ := json.Marshal(payload)

	fullURL := fmt.Sprintf("%s/content?token=%s",
		os.Getenv("BROWSERLESS_URL"),
		os.Getenv("BROWSERLESS_TOKEN"),
	)

	req, _ := http.NewRequestWithContext(ctx, "POST", fullURL, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("browserless error: %s", string(b))
	}

	data, err := io.ReadAll(resp.Body)
	return string(data), err
}
