package services

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/google/uuid"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/scheduler/models"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/queue"
	"gorm.io/gorm"
)

type FunctionService struct {
	db        *gorm.DB
	publisher queue.Publisher
}

func NewFunctionService(
	db *gorm.DB,
	publisher queue.Publisher,
) *FunctionService {
	return &FunctionService{
		db:        db,
		publisher: publisher,
	}
}

func (s *FunctionService) Run(function string, cpu string, memory string) error {
	if cpu == "" {
		cpu = "1"
	}

	if memory == "" {
		memory = "128mb"
	}

	functionName := fmt.Sprintf("%s-%s", function, uuid.NewString())

	_, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("docker run --rm -d --network=my-own-lambda --add-host=host.docker.internal:host-gateway --cpus %s --memory %s --name %s tiagorosadacosta123456/lambda-%s", cpu, memory, functionName, function)).Output()
	if err != nil {
		return err
	}
	return nil
}

func (s *FunctionService) ProcessLambdasScheduled() {
	var functions []models.Function
	s.db.Raw("SELECT *  FROM \"functions\" WHERE trigger = 'cron' and (last_execution + interval * interval '1 second') <= CURRENT_TIMESTAMP; ").Scan(&functions)
	var ids []int = []int{}
	for _, f := range functions {
		ids = append(ids, f.ID)
		go s.publisher.Publish(f, 1)
	}

	s.db.Model(models.Function{}).
		Where("trigger = 'cron'").
		Where("id in (?)", ids).
		Updates(models.Function{
			LastExecution: time.Now(),
		})

}
