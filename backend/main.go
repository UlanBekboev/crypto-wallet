package main

import (
	"backend/config"
	"backend/middleware"
	"backend/models"
	"backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"
	"log"
	"time"
	"os"
)

func main() {
	// Загрузка переменных окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Инициализация конфигурации приложения
	config.InitAppConfig()

	// Инициализация БД
	err = config.InitDB()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	// Миграции моделей
	err = config.DB.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Transaction{}, &models.PasswordReset{},)
	if err != nil {
		log.Fatal("Ошибка миграции:", err)
	}

	// Инициализация Redis Rate Limiter
	middleware.InitRateLimiter()

	// Запуск сервера
	r := gin.Default()
	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Добро пожаловать!"})
	})

	routes.AuthRoutes(r)
	routes.WalletRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err = r.Run(":" + port)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
