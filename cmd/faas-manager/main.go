package main

import (
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/configs"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/handler"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/services"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/queue"
	"gorm.io/gorm"
)

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	items = make(map[string]Item)
	mu    sync.Mutex
)

var db *gorm.DB

var publisher *queue.Publisher
var consumer *queue.Consumer

func init() {
	publisher = queue.NewPublisher("builder_docker_image")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db = configs.InitDB()
	app := fiber.New()

	functionService := services.NewFunctionService(db, *publisher)
	functionHandler := handler.NewFunctionHandler(
		*functionService,
	)

	app.Get("/functions/:id", functionHandler.FindById)
	app.Get("/functions", functionHandler.FindAll)
	app.Post("/functions", functionHandler.Deploy)

	app.Listen(":3000")
}
