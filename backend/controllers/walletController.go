package controllers

import (
	"backend/config"
	"backend/models"
	"gorm.io/gorm"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// Получить баланс кошелька
func GetWallet(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}

	var wallet models.Wallet
	if err := config.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Кошелёк не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": wallet.Balance})
}

// Перевод средств между кошельками
func Transfer(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}

	var input struct {
		ToUserID string `json:"to_user_id" binding:"required"`
		Amount   int64  `json:"amount" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	toID, err := uuid.Parse(input.ToUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный UUID получателя"})
		return
	}

	// Транзакция БД
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		var fromWallet, toWallet models.Wallet

		// Получаем кошельки
		if err := tx.Where("user_id = ?", userID).First(&fromWallet).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", toID).First(&toWallet).Error; err != nil {
			return err
		}

		// Проверка баланса
		if fromWallet.Balance < input.Amount {
			return errors.New("недостаточно средств")
		}

		// Обновляем балансы
		fromWallet.Balance -= input.Amount
		toWallet.Balance += input.Amount

		if err := tx.Save(&fromWallet).Error; err != nil {
			return err
		}
		if err := tx.Save(&toWallet).Error; err != nil {
			return err
		}

		// Записываем транзакции
		outTx := models.Transaction{
			ID:       uuid.New(),
			WalletID: fromWallet.ID,
			Amount:   input.Amount,
			Type:     "out",
		}
		inTx := models.Transaction{
			ID:       uuid.New(),
			WalletID: toWallet.ID,
			Amount:   input.Amount,
			Type:     "in",
		}
		if err := tx.Create(&outTx).Error; err != nil {
			return err
		}
		if err := tx.Create(&inTx).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Перевод успешно выполнен"})
}

// История транзакций
func TransactionHistory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}

	var wallet models.Wallet
	if err := config.DB.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Кошелёк не найден"})
		return
	}

	var transactions []models.Transaction
	if err := config.DB.Where("wallet_id = ?", wallet.ID).Order("created_at desc").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении истории"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
