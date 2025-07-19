// config/config.go
package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() error {
	dsn := "host=localhost user=postgres password=yourpassword dbname=wallet port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	log.Println("Подключение к базе данных прошло успешно!")
	return nil
}
