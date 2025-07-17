package utils

import (
	"github.com/gin-gonic/gin"
)

func SetAccessTokenCookie(c *gin.Context, token string) {
	c.SetCookie("access_token", token, 3600, "/", "", false, true)
}

func SetRefreshTokenCookie(c *gin.Context, token string) {
	c.SetCookie("refresh_token", token, 7*24*3600, "/", "", false, true)
}
