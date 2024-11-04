package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/log-collector/configs"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/log-collector/jobs"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := configs.Init()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	jobs.Init(db)

}
