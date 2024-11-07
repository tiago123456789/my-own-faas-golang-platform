package services

import (
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/tiago123456789/my-own-faas-golang-platform/internal/proxy/repositories"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/cache"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/queue"
)

type FunctionExecutorService struct {
	cache      cache.Cache
	publisher  queue.Publisher
	repository *repositories.FunctionRepository
}

func NewFunctionExecutorService(
	cache cache.Cache,
	publisher queue.Publisher,
	repository *repositories.FunctionRepository,
) *FunctionExecutorService {
	return &FunctionExecutorService{
		cache:      cache,
		publisher:  publisher,
		repository: repository,
	}
}

func (f *FunctionExecutorService) delete(function string) error {
	_, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("docker rm -f %s", function)).Output()
	if err != nil {
		return err
	}
	return nil
}

func (f *FunctionExecutorService) start(function string, cpu string, memory string) error {
	if cpu == "" {
		cpu = "1"
	}

	if memory == "" {
		memory = "128mb"
	}

	_, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("docker run --rm -d --network=my-own-lambda --add-host=host.docker.internal:host-gateway --cpus %s --memory %s --name %s tiagorosadacosta123456/lambda-%s", cpu, memory, function, function)).Output()
	if err != nil {
		return err
	}
	return nil
}

func (f *FunctionExecutorService) Stop(function string) error {
	functionCached, _ := f.cache.Get(function)
	if functionCached == "" {
		f.cache.Del(function)
		f.delete(function)
		err := f.delete(function)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *FunctionExecutorService) Run(function string) error {
	functionCached, _ := f.cache.Get(function)
	if functionCached == "" {
		functionReturned := f.repository.FindByName(function)
		if functionReturned.ID == 0 {
			return errors.New("Function not found")
		}

		err := f.start(function, functionReturned.Cpu, functionReturned.Memory)
		if err != nil {
			fmt.Printf("Error: %v", err)
			f.delete(function)
			f.start(function, functionReturned.Cpu, functionReturned.Memory)
		}

		go func() {
			f.cache.Set(function, true, 6*time.Second)
			f.publisher.PublishWithDelay(
				map[string]interface{}{
					"name": function,
				},
				2,
				60*time.Second,
			)
		}()
	} else {
		go f.cache.Set(function, true, 60*time.Second)
	}

	return nil
}
