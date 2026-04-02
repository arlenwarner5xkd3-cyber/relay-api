package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Item struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := fiber.New(fiber.Config{
		AppName:      "Relay API v0.1.0",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	api := app.Group("/api")
	api.Get("/items", listItems)
	api.Post("/items", createItem)
	api.Get("/items/:id", getItem)
	api.Delete("/items/:id", deleteItem)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down gracefully...")
	_ = app.Shutdown()
}

var items []Item

func listItems(c *fiber.Ctx) error {
	return c.JSON(items)
}

func createItem(c *fiber.Ctx) error {
	var item Item
	if err := c.BodyParser(&item); err != nil {
		return c.Status(400).JSON(ErrorResponse{Error: "bad_request", Message: "invalid JSON body"})
	}
	if item.Title == "" {
		return c.Status(400).JSON(ErrorResponse{Error: "validation", Message: "title is required"})
	}
	item.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	item.CreatedAt = time.Now().UTC()
	items = append(items, item)
	return c.Status(201).JSON(item)
}

func getItem(c *fiber.Ctx) error {
	id := c.Params("id")
	for _, item := range items {
		if item.ID == id {
			return c.JSON(item)
		}
	}
	return c.Status(404).JSON(ErrorResponse{Error: "not_found", Message: "item not found"})
}

func deleteItem(c *fiber.Ctx) error {
	id := c.Params("id")
	for i, item := range items {
		if item.ID == id {
			items = append(items[:i], items[i+1:]...)
			return c.SendStatus(204)
		}
	}
	return c.Status(404).JSON(ErrorResponse{Error: "not_found", Message: "item not found"})
}
