package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitDB() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл")
	}

	dbURL := os.Getenv("DATABASE_URL")
	DB, err = sqlx.Open("postgres", dbURL)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	log.Println("✅ Успешное подключение к PostgreSQL")
	return nil
}
