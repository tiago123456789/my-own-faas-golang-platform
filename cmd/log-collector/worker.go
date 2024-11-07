package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/log-collector/configs"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/log-collector/jobs"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/log-collector/repositories"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := configs.InitDB()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	logRepository := repositories.NewLogRepository(db)
	jobs.Init(logRepository)
}
