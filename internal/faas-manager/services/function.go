package services

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/models"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/repositories"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/types"
	"github.com/tiago123456789/my-own-faas-golang-platform/pkg/queue"
)

type FunctionService struct {
	publisher    queue.Publisher
	repositories repositories.FunctionRepository
}

func NewFunctionService(
	publisher queue.Publisher,
	repositories repositories.FunctionRepository,
) *FunctionService {
	return &FunctionService{
		publisher:    publisher,
		repositories: repositories,
	}
}

func (f *FunctionService) GetLogs(functionName string) []types.Log {
	return f.repositories.GetLogs(functionName)
}

func (f *FunctionService) FindById(id string) models.Function {
	return f.repositories.FindById(id)
}

func (f *FunctionService) FindAll() []models.Function {
	return f.repositories.FindAll()
}

func (f *FunctionService) isIntervalValid(interval string) error {
	regexExtractNumber := regexp.MustCompile(`([^0-9])+`)
	regexExtractNoNumerical := regexp.MustCompile(`([^a-z])+`)

	intervalNumber := regexExtractNumber.ReplaceAllString(interval, "")
	intervalNoNumber := regexExtractNoNumerical.ReplaceAllString(interval, "")

	if intervalNoNumber != "m" && intervalNoNumber != "h" {
		return errors.New("The interval format is wrong. The format allowed are: 1m(1 minute), 5m(5 minute), 5h(5 hours) and 1h(1 hour).")
	}

	intervalTypeInt, _ := strconv.Atoi(intervalNumber)
	if intervalTypeInt < 1 {
		return errors.New("The mininum time is 1 minute to lambda function using cron trigger.")
	}

	return nil
}

func (f *FunctionService) isRuntimeAllowed(runtime string) error {
	runtimesAllowed := map[string]bool{
		"golang:1.20":      true,
		"golang:1.19":      true,
		"golang:1.23":      true,
		"golang-cron:1.23": true,
		"golang-cron:1.20": true,
		"golang-cron:1.19": true,
	}

	if runtimesAllowed[runtime] == false {
		return errors.New("You provided the invalid runtime. The runtimes allowed are: golang:1.20, golang:1.19, golang:1.23, golang-cron:1.23, golang-cron:1.20 and golang-cron:1.19")
	}

	return nil
}

func (f *FunctionService) getIntervalInSeconds(interval string) int {
	regexExtractNumber := regexp.MustCompile(`([^0-9])+`)
	regexExtractNoNumerical := regexp.MustCompile(`([^a-z])+`)

	intervalNumber := regexExtractNumber.ReplaceAllString(interval, "")
	intervalNoNumber := regexExtractNoNumerical.ReplaceAllString(interval, "")

	actionToApply := map[string]func(value int) int{
		"m": func(value int) int {
			return value * 60
		},
		"h": func(value int) int {
			return value * (60 * 60)
		},
	}

	intervalTypeInt, _ := strconv.Atoi(intervalNumber)
	return actionToApply[intervalNoNumber](intervalTypeInt)

}

func (f *FunctionService) Deploy(newFunction types.NewFunction, lambdaPath string) (int, error) {
	intervalToSave := 0

	err := f.isIntervalValid(newFunction.Interval)
	if newFunction.Trigger == "cron" && err != nil {
		return 0, err
	}

	err = f.isRuntimeAllowed(newFunction.Runtime)
	if err != nil {
		return 0, err
	}

	if newFunction.Trigger == "cron" {
		intervalToSave = f.getIntervalInSeconds(newFunction.Interval)
	}

	function := models.Function{
		LambdaName:    newFunction.Name,
		Runtime:       newFunction.Runtime,
		LambdaPath:    lambdaPath,
		Cpu:           newFunction.Cpu,
		Memory:        newFunction.Memory,
		BuildProgress: "PENDENT",
		Trigger:       newFunction.Trigger,
		LastExecution: time.Now().Local().Add(5 * time.Minute),
		Interval:      intervalToSave,
	}

	functionReturned := f.repositories.FindByName(newFunction.Name)
	if functionReturned.ID == 0 {
		f.repositories.Create(&function)
	} else {
		function.ID = functionReturned.ID
		f.repositories.Update(function)
	}

	newFunction.ID = function.ID
	newFunction.LambdaPath = lambdaPath
	err = f.publisher.Publish(newFunction, 2)
	if err != nil {
		return 0, err
	}

	return function.ID, nil
}
