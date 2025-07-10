package store

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Import the Postgres dialect for GORM
	"github.com/rsys-speerzad/stackgen/pkg/models"
)

var db *gorm.DB

func InitDB() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		panic("DB_HOST environment variable is not set")
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		panic("DB_PORT environment variable is not set")
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		panic("DB_USER environment variable is not set")
	}
	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		panic("DB_PASS environment variable is not set")
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		panic("DB_NAME environment variable is not set")
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	var err error
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
}

func GetDB() *gorm.DB {
	if db == nil {
		panic("database connection is not initialized")
	}
	return db
}

func CloseDB() {
	if db != nil {
		if err := db.Close(); err != nil {
			fmt.Printf("failed to close database connection: %v\n", err)
		}
		db = nil
	} else {
		fmt.Println("database connection is already closed or not initialized")
	}
}

// AutoMigrate runs the database migrations.
func AutoMigrate() error {
	if db == nil {
		InitDB()
	}

	if err := db.AutoMigrate(
		&models.Event{},
		&models.User{},
		&models.EventSlot{},
		&models.UserAvailability{},
	).Error; err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
