package routes

import (
	"backend/controllers"
	"backend/middleware"
	"github.com/gin-gonic/gin"
)

func WalletRoutes(r *gin.Engine) {
	auth := r.Group("/api/auth/wallet")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/connect", controllers.ConnectWallet)
		auth.POST("/", controllers.CreateWallet)
		auth.GET("/", controllers.GetWallet)
	}
}
