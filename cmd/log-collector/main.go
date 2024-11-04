package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/log-collector/handler"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/queue"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logPublisher := queue.NewPublisher("logs")

	app := fiber.New()

	logHandler := handler.NewLogHandler(*logPublisher)

	app.Post("/", logHandler.Register)

	app.Listen(":5050")
}
