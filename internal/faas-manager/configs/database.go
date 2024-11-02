package configs

import (
	"log"
	"os"

	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dbURL := os.Getenv("DB_URL")

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.Function{})

	return db
}
