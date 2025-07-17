package controllers

import (
	"backend/config"
	"backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ConnectWallet(c *gin.Context) {
	var input struct {
		Address   string `json:"address"`
		Signature string `json:"signature"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат"})
		return
	}

	userID := c.GetInt("user_id")
	message := "Подтверждение владения кошельком"

	err := utils.VerifySignature(input.Address, input.Signature, message)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Обновляем wallet_address в users
	_, err = config.DB.Exec("UPDATE users SET wallet_address=$1 WHERE id=$2", input.Address, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении адреса"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Кошелёк успешно привязан"})
}
