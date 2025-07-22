// jwt.go
package utils

import (
	"backend/config"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func GenerateAccessToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWT_SECRET))
}

func GenerateRefreshToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWT_SECRET))
}

func GenerateTokens(userID uuid.UUID) (Tokens, error) {
	access, err := GenerateAccessToken(userID)
	if err != nil {
		return Tokens{}, err
	}
	refresh, err := GenerateRefreshToken(userID)
	if err != nil {
		return Tokens{}, err
	}
	return Tokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func SetAuthCookies(c *gin.Context, tokens Tokens) {
	c.SetCookie("access_token", tokens.AccessToken, 60*15, "/", "", false, true)
	c.SetCookie("refresh_token", tokens.RefreshToken, 60*60*24*7, "/", "", false, true)
}

func ValidateToken(tokenString string, secret string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	return claims.Subject, nil
}

func ParseRefreshToken(tokenString string) (string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_SECRET), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	if claims["type"] != "refresh" {
		return "", errors.New("invalid token type")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("invalid user_id")
	}

	return userID, nil
}

