package main

import (
	"backend/config"
	"backend/routes"
	"backend/middleware"
	"backend/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"gorm.io/gorm"
	"log"
)

func main() {
	err := config.InitDB()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	err = config.DB.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Transaction{})
	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	middleware.InitRateLimiter()

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Добро пожаловать!"})
	})

	routes.AuthRoutes(r)
	routes.WalletRoutes(r)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
