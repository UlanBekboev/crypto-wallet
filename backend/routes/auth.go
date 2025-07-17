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
		auth.GET("/profile", middleware.AuthMiddleware(), controllers.GetProfile)
		auth.PUT("/profile", middleware.AuthMiddleware(), controllers.UpdateProfile)
		auth.POST("/forgot-password", controllers.ForgotPassword)
		auth.POST("/reset-password", controllers.ResetPassword)
		auth.POST("/logout", controllers.Logout)
	}
}
