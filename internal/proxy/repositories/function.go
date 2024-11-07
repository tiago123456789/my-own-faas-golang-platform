package repositories

import (
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/proxy/models"
	"gorm.io/gorm"
)

type FunctionRepository struct {
	db *gorm.DB
}

func NewFunctionRepository(
	db *gorm.DB,
) *FunctionRepository {
	return &FunctionRepository{
		db: db,
	}
}

func (r *FunctionRepository) FindByName(name string) models.Function {
	var function models.Function
	r.db.First(&function, "lambda_name = ?", name)
	return function
}
