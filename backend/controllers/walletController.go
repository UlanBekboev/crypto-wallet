package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"backend/models"
	"backend/config"
)

// POST /wallet
func CreateWallet(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	id := uuid.New()
	_, err := config.DB.Exec(`
		INSERT INTO wallets (id, user_id, balance, created_at)
		VALUES ($1, $2, 0, NOW())
	`, id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "wallet creation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet_id": id})
}

// GET /wallet
func GetWallet(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	var wallet models.Wallet
	err := config.DB.Get(&wallet, `SELECT * FROM wallets WHERE user_id = $1`, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}