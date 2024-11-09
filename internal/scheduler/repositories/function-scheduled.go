package repositories

import (
	"time"

	"github.com/tiago123456789/my-own-faas-golang-platform/internal/scheduler/models"
	"gorm.io/gorm"
)

type FunctionScheduledRepository struct {
	db *gorm.DB
}

func NewFunctionScheduledRepository(db *gorm.DB) *FunctionScheduledRepository {
	return &FunctionScheduledRepository{
		db: db,
	}
}

func (r *FunctionScheduledRepository) GetFunctionsNeedsToProcess() []models.Function {
	var functions []models.Function
	r.db.Raw("SELECT *  FROM \"functions\" WHERE build_progress = 'DONE' and trigger = 'cron' and (last_execution + interval * interval '1 second') <= CURRENT_TIMESTAMP; ").Scan(&functions)
	return functions
}

func (r *FunctionScheduledRepository) UpdateLastExecutionByIds(ids []int) {
	r.db.Model(models.Function{}).
		Where("trigger = 'cron'").
		Where("id in (?)", ids).
		Updates(models.Function{
			LastExecution: time.Now(),
		})
}
