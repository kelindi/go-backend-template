package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

// Database connection
func Connect() (*gorm.DB, error) {
	dsn := "user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" host=" + os.Getenv("DB_HOST") +
		" dbname=" + os.Getenv("DB_DATABASE") +
		" sslmode=" + os.Getenv("DB_SSLMODE")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	//conection pooling settings
	sqlDB.SetMaxIdleConns(20)

	sqlDB.SetMaxOpenConns(90)

	sqlDB.SetConnMaxLifetime(time.Hour)

	// Attempt to enable the uuid-ossp extension
	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		return nil, fmt.Errorf("failed to enable uuid-ossp extension: %w", err)
	}

	if err := createSchema(db); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return db, nil
}

func createSchema(db *gorm.DB) error {
	models := []interface{}{
		//add models here
		&User{},
	}

	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			return err
		}
	}

	return nil
}
