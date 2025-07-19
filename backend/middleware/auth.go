package middleware

import (
	"backend/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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

// Middleware для проверки JWT access токена
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("access_token")
		if err != nil || tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Отсутствует токен авторизации",
			})
			return
		}

		claims, err := utils.ValidateAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Неверный или просроченный токен",
			})
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "user_id не строка",
			})
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "user_id невалидный UUID",
			})
			return
		}

		// Передаём user_id в контекст
		c.Set("user_id", userID)
		c.Next()
	}
}
