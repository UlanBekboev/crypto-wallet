package routes

import (
	"backend/controllers"
	"backend/middleware"
	"github.com/gin-gonic/gin"
)

func WalletRoutes(r *gin.Engine) {
	auth := r.Group("/api/wallet")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.GET("/", controllers.GetWallet)
		auth.POST("/transfer", controllers.Transfer)
		auth.GET("/history", controllers.TransactionHistory)
	}
}
