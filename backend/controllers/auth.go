package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-playground/validator/v10"
	"strings"
	"fmt"
)

var validate = validator.New()

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))

	var user models.User
	err := config.DB.Get(&user, "SELECT id, email, password FROM users WHERE email=$1", email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный пароль"})
		return
	}

	// Generate JWT tokens
	accessToken, _ := utils.GenerateAccessToken(user.ID)
	refreshToken, _ := utils.GenerateRefreshToken(user.ID)

	// Устанавливаем токены в cookie
	c.SetCookie("access_token", accessToken, 3600, "/", "localhost", false, true)
	c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Успешный вход",
	})
}


func Register(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	if err := utils.Validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors": utils.FormatValidationErrors(err),
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка хеширования"})
		return
	}

	var existing models.User
	err = config.DB.Get(&existing, "SELECT * FROM users WHERE email=$1", input.Email)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email уже зарегистрирован"})
		return
	}

	_, err = config.DB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", input.Email, hashedPassword)
if err != nil {
    if strings.Contains(err.Error(), "duplicate key value") && strings.Contains(err.Error(), "users_email_key") {
        c.JSON(http.StatusConflict, gin.H{"error": "Этот email уже зарегистрирован"})
    } else {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Ошибка записи в БД",
            "details": err.Error(),
        })
    }
    return
}

	c.JSON(http.StatusCreated, gin.H{"message": "Успешная регистрация"})
}

func RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh токен отсутствует"})
		return
	}

	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Недействительный refresh токен"})
		return
	}

	userID := int(claims["user_id"].(float64))
	newAccessToken, _ := utils.GenerateAccessToken(userID)

	c.SetCookie("access_token", newAccessToken, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Токен обновлён",
	})
}


func Me(c *gin.Context) {
	userID := c.GetInt("user_id")

	var user models.User
	err := config.DB.Get(&user, "SELECT id, email, name, created_at FROM users WHERE id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Пользователь не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name": func() string {
				if user.Name.Valid {
					return user.Name.String
				}
				return ""
			}(),
		},
	})
}

func ChangePassword(c *gin.Context) {
	type PasswordInput struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	var input PasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	userID := c.GetInt("user_id")

	var hashedPassword string
	err := config.DB.Get(&hashedPassword, "SELECT password FROM users WHERE id=$1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Пользователь не найден"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Старый пароль неверен"})
		return
	}

	newHashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось зашифровать новый пароль"})
		return
	}

	_, err = config.DB.Exec("UPDATE users SET password=$1 WHERE id=$2", string(newHashed), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении пароля"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пароль успешно изменён"})
}

func GetProfile(c *gin.Context) {
	userID := c.GetInt("user_id")
	fmt.Println("userID из токена:", userID)
	var user models.User

	err := config.DB.Get(&user, "SELECT id, email, name, created_at FROM users WHERE id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка получения профиля",
			"details": err.Error(),
		})
		
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"name":       func() string {
				if user.Name.Valid {
					return user.Name.String
				}
				return ""
			}(),
		"created_at": user.CreatedAt,
	})
}

func UpdateProfile(c *gin.Context) {
	userID := c.GetInt("user_id")

	type Input struct {
		Name string `json:"name"`
	}

	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	_, err := config.DB.Exec("UPDATE users SET name=$1 WHERE id=$2", input.Name, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить профиль"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Профиль обновлён"})
}

func Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Выход выполнен"})
}