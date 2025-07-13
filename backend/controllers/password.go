package controllers

import (
	"backend/config"
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

	// Сгенерировать токен
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	token := hex.EncodeToString(b)

	// Сохранить токен в БД
	_, err := config.DB.Exec(`
		INSERT INTO password_resets (email, token, expires_at)
		VALUES ($1, $2, $3)
	`, input.Email, token, time.Now().Add(15*time.Minute))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать токен"})
		return
	}

	resetLink := fmt.Sprintf("http://localhost:3000/reset-password?token=%s&email=%s", token, input.Email)

	err = utils.SendEmail(input.Email, "Сброс пароля", "Ссылка для сброса пароля: "+resetLink)
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
        "error":   "Не удалось отправить email",
        "details": err.Error(), // добавим вывод самой ошибки
    })
    return
}

	c.JSON(http.StatusOK, gin.H{"message": "Письмо отправлено"})
}

func ResetPassword(c *gin.Context) {
	var input struct {
		Token       string `json:"token"`
		NewPassword string `json:"newPassword"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	// 1. Найти токен в password_resets
	var email string
	var expiresAt time.Time
	err := config.DB.QueryRow(`
		SELECT email, expires_at FROM password_resets WHERE token=$1
	`, input.Token).Scan(&email, &expiresAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный токен"})
		return
	}

	// 2. Проверить срок действия
	if time.Now().After(expiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Токен истёк"})
		return
	}

	// 3. Хешируем новый пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хешировании"})
		return
	}

	// 4. Обновляем пароль
	_, err = config.DB.Exec(`UPDATE users SET password=$1 WHERE email=$2`, string(hashedPassword), email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить пароль"})
		return
	}

	// 5. Удаляем использованный токен
	config.DB.Exec(`DELETE FROM password_resets WHERE token=$1`, input.Token)

	c.JSON(http.StatusOK, gin.H{"message": "Пароль успешно сброшен"})
}

