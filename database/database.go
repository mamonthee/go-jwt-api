package database

import (
	"go-jwt-api/models"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error

	// Set the log level to silent (no logs will be shown)
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: 200 * time.Millisecond, // Slow SQL threshold
			LogLevel:      logger.Silent,          // Silent logging
			Colorful:      false,                  // No color for log output
		},
	)

	DB, err = gorm.Open(sqlite.Open("content_management.db"), &gorm.Config{Logger: gormLogger})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate models
	err = DB.AutoMigrate(&models.Author{}, &models.Article{})
	if err != nil {
		log.Fatalf("Failed to migrate models: %v", err)
	}
}
