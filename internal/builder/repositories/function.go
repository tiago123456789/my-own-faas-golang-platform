package repositories

import (
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

func (r *FunctionRepository) UpdateProcess(id interface{}, status string) {
	r.db.Exec(
		"UPDATE functions SET build_progress = ? WHERE id = ?",
		status,
		id,
	)
}
