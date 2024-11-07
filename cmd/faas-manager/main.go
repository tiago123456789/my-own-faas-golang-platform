package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/configs"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/handler"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/repositories"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/services"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/queue"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := configs.InitDB()
	esDB, err := configs.InitES()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	publisher := queue.NewPublisher("builder_docker_image")

	functionRepository := repositories.NewFunctionRepository(
		db, esDB,
	)
	functionService := services.NewFunctionService(*publisher, *functionRepository)
	functionHandler := handler.NewFunctionHandler(
		*functionService,
	)

	app.Get("/functions/:id/logs", functionHandler.GetLogs)
	app.Get("/functions/:id", functionHandler.FindById)
	app.Get("/functions", functionHandler.FindAll)
	app.Post("/functions", functionHandler.Deploy)

	app.Listen(":3000")
}
