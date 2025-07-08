package main

import (
	"backend/config"
	"backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Backend работает!"})
	})

	routes.AuthRoutes(router)

	router.Run(":" + config.GetEnv("PORT"))
}
