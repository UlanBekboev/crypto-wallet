package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"golang.org/x/crypto/bcrypt"
)

func ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	// Генерация токена
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	token := hex.EncodeToString(b)

	expiration := time.Now().Add(15 * time.Minute)

	// Сохраняем токен в БД
	reset := models.PasswordReset{
		Email:     input.Email,
		Token:     token,
		ExpiresAt: expiration,
	}
	config.DB.Save(&reset)

	resetLink := fmt.Sprintf("http://localhost:3000/reset-password?token=%s&email=%s", token, input.Email)

	err := utils.SendEmail(input.Email, "Сброс пароля", "Ссылка для сброса пароля: "+resetLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось отправить email", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Письмо отправлено"})
}


func ResetPassword(c *gin.Context) {
	var input struct {
		Token       string `json:"token"`
		Email       string `json:"email"`
		NewPassword string `json:"newPassword"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	var reset models.PasswordReset
	if err := config.DB.Where("email = ? AND token = ?", input.Email, input.Token).First(&reset).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный токен или email"})
		return
	}

	if time.Now().After(reset.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Токен истёк"})
		return
	}

	// Найти пользователя
	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь не найден"})
		return
	}

	// Хешируем новый пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хешировании"})
		return
	}

	user.Password = string(hashedPassword)
	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить пароль"})
		return
	}

	// Удаляем использованный токен
	config.DB.Delete(&reset)

	c.JSON(http.StatusOK, gin.H{"message": "Пароль успешно обновлён"})
}

