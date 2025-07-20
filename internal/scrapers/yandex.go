package scrapers

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"browserless-image-api/internal/fetcher"

	"github.com/PuerkitoBio/goquery"
)

type YandexScraper struct{}

func (y YandexScraper) SearchImages(ctx context.Context, query string, max int) ([]string, error) {
	searchURL := fmt.Sprintf("https://yandex.com/images/search?text=%s", url.QueryEscape(query))
	html, err := fetcher.FetchContentStealth(ctx, searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch search page: %w", err)
	}

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	hrefs := make([]string, 0, max)

	doc.Find("div.ImagesContentImage > a.ImagesContentImage-Cover").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if len(hrefs) >= max {
			return false
		}
		if href, ok := s.Attr("href"); ok {
			hrefs = append(hrefs, "https://yandex.com"+href)
		}
		return true
	})

	var images []string

	for _, thumb := range hrefs {
		parsed, err := url.Parse(thumb)
		if err != nil {
			continue
		}

		q := parsed.Query()
		imgURL := q.Get("img_url")
		if imgURL == "" {
			continue
		}

		images = append(images, imgURL)
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("no image URLs extracted")
	}

	return images, nil
}
