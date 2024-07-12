// Роуты для управления кошелькам в базе данных
package handlers

import (
	"dao_vote/back-end/repository"
	"dao_vote/back-end/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// WalletStrength представляет структуру для добавления силы голоса кошелька
type WalletStrength struct {
	WalletAddress string `json:"wallet_address" binding:"required"`
	VotePower     int    `json:"vote_power" binding:"required"`
}

// AddWalletHandler добавляет новый адрес кошелька и силу голоса в базу данных
func AddWalletHandler(c *gin.Context) {
	utils.HandleRequest(c, func(c *gin.Context) error {
		var wallet WalletStrength
		if err := c.ShouldBindJSON(&wallet); err != nil {
			logrus.Errorf("Invalid request body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return err
		}

		// Проверяем, существует ли уже запись с таким wallet_address
		if _, err := repository.GetVoteStrength(wallet.WalletAddress); err == nil {
			logrus.Errorf("Wallet address already exists: %v", wallet.WalletAddress)
			c.JSON(http.StatusConflict, gin.H{"error": "Wallet address already exists"})
			return err
		}

		if err := repository.AddWalletStrength(wallet.WalletAddress, wallet.VotePower); err != nil {
			logrus.Errorf("Failed to add wallet: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add wallet"})
			return err
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Wallet added successfully"})
		logrus.Info("AddWalletHandler completed successfully")
		return nil
	})
}

// DeleteWalletHandler удаляет адрес кошелька и силу голоса из базы данных
func DeleteWalletHandler(c *gin.Context) {
	utils.HandleRequest(c, func(c *gin.Context) error {
		walletAddress := c.Param("wallet_address")
		if err := repository.DeleteWalletStrength(walletAddress); err != nil {
			logrus.Errorf("Failed to delete wallet: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete wallet"})
			return err
		}

		c.JSON(http.StatusOK, gin.H{"message": "Wallet deleted successfully"})
		logrus.Info("DeleteWalletHandler completed successfully")
		return nil
	})
}
