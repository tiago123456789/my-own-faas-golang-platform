package main

import (
	"log"

	"github.com/joho/godotenv"
	job "github.com/tiago123456789/my-own-faas-golang-platform/internal/builder/jobs"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	job.Init()
}
