package main

import (
	"backend/config"
	"backend/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	err := config.InitDB()
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Добро пожаловать!"})
	})

	routes.AuthRoutes(r)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

/* func main() {
	config.LoadEnv()
	config.ConnectDB()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Backend работает!"})
	})

	routes.AuthRoutes(router)

	router.Run(":" + config.GetEnv("PORT"))
} */

