package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/proxy/configs"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/proxy/handler"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/proxy/jobs"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/proxy/services"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/cache"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/queue"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := configs.InitDB()
	app := fiber.New()

	ctx := context.Background()
	cache, err := cache.New(ctx)
	defer cache.Close()
	if err != nil {
		log.Fatalln("Redis connection was refused")
	}

	publisher := queue.NewPublisher(
		"delete_function_with_expire",
	)

	functionExecutor := services.NewFunctionExecutorService(
		*cache, *publisher, db,
	)

	jobs.Init(*functionExecutor)

	functionHandler := handler.NewFunctionHandler(
		*functionExecutor,
	)

	app.Use("/:function/*", functionHandler.Execute)

	app.Listen(":8080")
}
