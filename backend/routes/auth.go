package routes

import (
	"backend/controllers"
	"backend/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.POST("/refresh", controllers.RefreshToken)
		auth.GET("/me", middleware.AuthMiddleware(), controllers.Me)
		auth.POST("/change-password", middleware.AuthMiddleware(), controllers.ChangePassword)
		//auth.POST("/forgot-password", controllers.ForgotPassword)
		//auth.POST("/reset-password", controllers.ResetPassword)
		//auth.GET("/logout", controllers.Logout)
	}
}

