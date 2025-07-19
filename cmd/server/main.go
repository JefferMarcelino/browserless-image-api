package main

import (
	"browserless-image-api/config"
	"browserless-image-api/internal/scrapers"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()
	app := fiber.New()

	scraper := scrapers.YandexScraper{}

	app.Get("/image", func(c *fiber.Ctx) error {
		q := c.Query("q")
		if q == "" {
			return c.Status(400).SendString("missing q param")
		}

		max, _ := strconv.Atoi(c.Query("max", "3"))
		if max <= 0 {
			max = 5
		}

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		images, err := scraper.SearchImages(ctx, q, max)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(fiber.Map{
			"query":  q,
			"images": images,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
