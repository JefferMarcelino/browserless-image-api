package scrapers

import (
	"context"
)

type ImageScraper interface {
	SearchImages(ctx context.Context, query string, max int) ([]string, error)
}
