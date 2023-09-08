package repository

import (
	"fmt"
	"github.com/brcodingdev/chat-app/service/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	db *gorm.DB
)

// Connect ...
func Connect() (*gorm.DB, error) {
	dbUserName := os.Getenv("DB_USERNAME")
	dbUserPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")

	metricsLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUserName, dbUserPassword, dbName, dbPort,
	)
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: metricsLogger})

	if err != nil {
		return nil, err
	}
	db = d
	return db, nil
}

// MigrateDB ...
func MigrateDB() error {
	return db.AutoMigrate(model.Tables...)
}

// ErrorCheck ...
func ErrorCheck(db *gorm.DB) error {
	if err := db.Error; err != nil {
		return err
	}

	return nil
}
