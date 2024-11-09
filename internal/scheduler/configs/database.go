package configs

import (
	"log"
	"os"

	"github.com/tiago123456789/my-own-faas-golang-platform/internal/scheduler/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	dbURL := os.Getenv("DB_URL")

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.Function{})

	return db
}
