package fetcher

import (
	"context"
	"errors"
	"sync/atomic"

	g "github.com/serpapi/google-search-results-golang"
)

type ApiRotator struct {
	keys []string
	idx  uint32
}

type ImageResult struct {
	Original  string `json:"original"`
	Thumbnail string `json:"thumbnail"`
	Source    string `json:"source"`
	Link      string `json:"link"`
}

func NewApiRotator(keys []string) *ApiRotator {
	return &ApiRotator{keys: keys}
}

func (r *ApiRotator) Next() string {
	i := atomic.AddUint32(&r.idx, 1)
	return r.keys[int(i)%len(r.keys)]
}

func SearchImages(ctx context.Context, rot *ApiRotator, query string, max int) ([]ImageResult, error) {
	params := map[string]string{
		"q":             query,
		"engine":        "google_images",
		"google_domain": "google.co.mz",
		"gl":            "mz",
		"hl":            "pt",
		"device":        "desktop",
	}

	var resp map[string]any
	var err error

	for range rot.keys {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		apiKey := rot.Next()
		client := g.NewGoogleSearch(params, apiKey)

		resp, err = client.GetJSON()
		if err != nil {
			continue
		}
		break
	}

	if resp == nil {
		return nil, errors.New("all API keys failed")
	}

	raw, ok := resp["images_results"].([]any)
	if !ok {
		return nil, errors.New("missing images_results in response")
	}

	results := []ImageResult{}
	for _, item := range raw {
		if len(results) >= max {
			break
		}
		m := item.(map[string]any)
		results = append(results, ImageResult{
			Original:  m["original"].(string),
			Thumbnail: m["thumbnail"].(string),
			Source:    m["source"].(string),
			Link:      m["link"].(string),
		})
	}

	return results, nil
}
