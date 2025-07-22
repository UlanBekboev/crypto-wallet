package middleware

import (
	"net/http"
	"time"

	"backend/config"
	"backend/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/go-redis/redis_rate/v10"
	"github.com/golang-jwt/jwt/v5"
)

var (
	RDB     *redis.Client
	limiter *redis_rate.Limiter
)

// Инициализация Redis клиента и лимитера
func InitRateLimiter() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // или вынеси в .env
	})
	limiter = redis_rate.NewLimiter(RDB)
}

// Middleware для ограничения частоты запросов
func RateLimitMiddleware(limit int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// Используем контекст запроса для отмены при завершении запроса
		res, err := limiter.Allow(c.Request.Context(), ip, redis_rate.Limit{
			Rate:   limit,
			Period: duration,
			Burst:  1,
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка лимитера",
			})
			return
		}

		if res.Allowed == 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Слишком много запросов, попробуйте позже",
			})
			return
		}

		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 🔧 1. Попытка получить токен из cookie
		token, err := c.Cookie("access_token")
		if err != nil || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization cookie missing"})
			return
		}

		// 🔒 2. Проверка токена
		claims := jwt.MapClaims{}
		parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWT_SECRET), nil
		})
		if err != nil || !parsedToken.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// 👤 3. Извлечение ID пользователя
		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		var user models.User
		if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		// ✅ Устанавливаем пользователя в контекст
		c.Set("user", user)
		c.Next()
	}
}
