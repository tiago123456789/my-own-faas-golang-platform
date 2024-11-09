package services

import (
	"fmt"
	"os/exec"

	"github.com/google/uuid"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/scheduler/repositories"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/queue"
)

type FunctionService struct {
	repository *repositories.FunctionScheduledRepository
	publisher  queue.Publisher
}

func NewFunctionService(
	repository *repositories.FunctionScheduledRepository,
	publisher queue.Publisher,
) *FunctionService {
	return &FunctionService{
		repository: repository,
		publisher:  publisher,
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
	functions := s.repository.GetFunctionsNeedsToProcess()
	var ids []int = []int{}
	for _, f := range functions {
		ids = append(ids, f.ID)
		go s.publisher.Publish(f, 1)
	}

	s.repository.UpdateLastExecutionByIds(ids)
}
