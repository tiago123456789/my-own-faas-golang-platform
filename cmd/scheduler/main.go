package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/scheduler/configs"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/scheduler/cron"
	job "github.com/tiago123456789/my-own-faas-golang-platform/internal/scheduler/jobs"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/scheduler/repositories"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/scheduler/services"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/queue"
)

func main() {
	enableScheduler := flag.Bool("enable_scheduler", false, "a bool")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := configs.InitDB()

	publisher := queue.NewPublisher("lambda_executions_triggered")
	functionSchedulerRepository := repositories.NewFunctionScheduledRepository(
		db,
	)
	functionService := services.NewFunctionService(
		functionSchedulerRepository, *publisher,
	)

	if *enableScheduler == true {

		cron.NewCron(
			10*time.Second,
			functionService.ProcessLambdasScheduled,
		).Start()

		fmt.Println("Loaded the scheduler")
	}

	job.Init(functionService)
}
