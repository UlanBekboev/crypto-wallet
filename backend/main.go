package main

import (
	"backend/config"
	"backend/routes"
	"backend/middleware"
	"github.com/gin-gonic/gin"
	"backend/utils"
	"github.com/gin-contrib/cors"
	"log"
)

func main() {
	err := config.InitDB()
	utils.InitValidator()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	middleware.InitRateLimiter()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Добро пожаловать!"})
	})

	routes.AuthRoutes(r)
	routes.WalletRoutes(r)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}

	r.Run(":8080")
}

