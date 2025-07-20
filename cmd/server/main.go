package main

import (
	"context"
	"fmt"
	"image-api/config"
	"image-api/internal/fetcher"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()
	app := fiber.New()

	rot := fetcher.NewApiRotator([]string{
		os.Getenv("SERP_KEY1"),
	})

	app.Get("/image", func(c *fiber.Ctx) error {
		q := c.Query("q")
		if q == "" {
			return c.Status(400).SendString("missing q param")
		}

		max, _ := strconv.Atoi(c.Query("max", "3"))
		if max <= 0 {
			max = 5
		}

		ctx := context.Background()

		imgs, err := fetcher.SearchImages(ctx, rot, q, max)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"query":  q,
			"images": imgs,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
