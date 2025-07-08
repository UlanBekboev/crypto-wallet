package config

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectDB() {
	var err error
	DB, err = sqlx.Connect("postgres", GetEnv("DB_URL"))
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	log.Println("✅ Успешное подключение к PostgreSQL")
}
