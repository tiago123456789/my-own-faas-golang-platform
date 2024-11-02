package services

import (
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/models"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/types"
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

func (f *FunctionService) FindAll() []models.Function {
	var functions []models.Function
	f.db.Find(&functions)
	return functions
}

func (f *FunctionService) Deploy(newFunction types.NewFunction, lambdaPath string) error {
	function := models.Function{
		LambdaName:    newFunction.Name,
		Runtime:       newFunction.Runtime,
		LambdaPath:    lambdaPath,
		Cpu:           newFunction.Cpu,
		Memory:        newFunction.Memory,
		BuildProgress: "PENDENT",
	}

	var functionReturned models.Function
	f.db.First(&functionReturned, "lambda_name = ?", newFunction.Name)
	if functionReturned.ID == 0 {
		f.db.Create(&function)
	} else {
		function.ID = functionReturned.ID
		f.db.Updates(function)
	}

	newFunction.ID = function.ID
	newFunction.LambdaPath = lambdaPath
	err := f.publisher.Publish(newFunction, 2)
	if err != nil {
		return err
	}

	return nil
}

func (f *FunctionService) UpdateProgress(id int, status string) error {
	var functionReturned models.Function
	f.db.First(&functionReturned, "id = ?", id)
	if functionReturned.ID != 0 {
		functionToUpdate := models.Function{
			ID:            id,
			BuildProgress: status,
		}
		f.db.Updates(functionToUpdate)
	}

	return nil
}
