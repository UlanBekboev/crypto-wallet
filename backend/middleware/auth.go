package middleware

import (
	"backend/utils"
	"net/http"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/go-redis/redis_rate/v10"
	"time"
)

var (
	RDB     *redis.Client
	limiter *redis_rate.Limiter
)

func InitRateLimiter() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", 
	})
	limiter = redis_rate.NewLimiter(RDB)
}

func RateLimitMiddleware(limit int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		res, err := limiter.Allow(context.Background(), ip, redis_rate.Limit{
			Rate:   limit,
			Period: duration,
			Burst:  1,
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ошибка лимитера"})
			return
		}

		if res.Allowed == 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Слишком много запросов, попробуйте позже"})
			return
		}

		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("access_token")
		if err != nil || tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Отсутствует токен авторизации"})
			return
		}

		claims, err := utils.ValidateAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
			return
		}

		// Сохраняем user_id в контексте запроса
		c.Set("user_id", int(claims["user_id"].(float64)))
		c.Next()
	}
}
