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
	var wallet WalletStrength
	if err := c.ShouldBindJSON(&wallet); err != nil {
		utils.JSONResponse(c, http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		logrus.Errorf("Invalid request body: %v", err)
		return
	}

	err := repository.AddWalletStrength(wallet.WalletAddress, wallet.VotePower)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": "Failed to add wallet"})
		logrus.Errorf("Failed to add wallet: %v", err)
		return
	}

	utils.JSONResponse(c, http.StatusCreated, gin.H{"message": "Wallet added successfully"})
	logrus.Info("AddWalletHandler completed successfully")
}

// DeleteWalletHandler удаляет адрес кошелька и силу голоса из базы данных
func DeleteWalletHandler(c *gin.Context) {
	walletAddress := c.Param("wallet_address")

	err := repository.DeleteWalletStrength(walletAddress)
	if err != nil {
		utils.JSONResponse(c, http.StatusInternalServerError, gin.H{"error": "Failed to delete wallet"})
		logrus.Errorf("Failed to delete wallet: %v", err)
		return
	}

	utils.JSONResponse(c, http.StatusOK, gin.H{"message": "Wallet deleted successfully"})
	logrus.Info("DeleteWalletHandler completed successfully")
}
