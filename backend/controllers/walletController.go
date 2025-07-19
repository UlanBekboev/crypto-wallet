package controllers

import (
	"backend/config"
	"backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func CreateWallet(c *gin.Context) {
	userCtx, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}
	user := userCtx.(models.User)

	var input struct {
		Address string `json:"address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wallet := models.Wallet{
		ID:            uuid.New(),
		UserID:        user.ID,
		WalletAddress: input.Address,
	}

	if err := config.DB.Create(&wallet).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании кошелька"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Кошелёк успешно создан", "wallet": wallet})
}

func GetWallet(c *gin.Context) {
	userCtx, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}
	user := userCtx.(models.User)

	var wallet models.Wallet
	if err := config.DB.Where("user_id = ?", user.ID).First(&wallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Кошелёк не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": wallet})
}
