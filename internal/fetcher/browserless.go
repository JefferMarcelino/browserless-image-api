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

type bqlResp struct {
	Data struct {
		HTML struct {
			HTML string `json:"html"`
		} `json:"html"`
		Verify struct {
			Found  bool `json:"found"`
			Solved bool `json:"solved"`
			Time   int  `json:"time"`
		} `json:"verify"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func FetchContentStealth(ctx context.Context, pageURL string) (string, error) {
	payload := map[string]any{
		"query": fmt.Sprintf(`
      mutation FetchAndBypass {
        goto(url: "%s", waitUntil: networkIdle) { status }
        html { html }
      }
    `, pageURL),
		"variables":     map[string]any{},
		"operationName": "FetchAndBypass",
	}

	bodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf(
			"%s/chrome/bql?token=%s&proxy=residential&proxySticky=true&proxyCountry=mz&humanlike=true&blockConsentModals=true",
			os.Getenv("BROWSERLESS_URL"),
			os.Getenv("BROWSERLESS_TOKEN"),
		),
		bytes.NewReader(bodyBytes),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)

	var r bqlResp
	if err := json.Unmarshal(b, &r); err != nil {
		return "", fmt.Errorf("JSON unmarshal failed: %w", err)
	}
	if len(r.Errors) > 0 {
		return "", fmt.Errorf("BQL error: %s", r.Errors[0].Message)
	}
	if r.Data.Verify.Found && !r.Data.Verify.Solved {
		return "", fmt.Errorf("cloudflare challenge detected but not solved")
	}

	return r.Data.HTML.HTML, nil
}
