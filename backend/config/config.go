package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"log"
)

var (
	DB         *gorm.DB
	JWT_SECRET string
)

func InitDB() error {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func InitAppConfig() {
	JWT_SECRET = os.Getenv("JWT_SECRET")
	if JWT_SECRET == "" {
		log.Fatal("JWT_SECRET not set in environment variables")
	}
}