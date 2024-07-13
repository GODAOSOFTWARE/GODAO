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
	logrus.Info("AddWalletHandler called")

	// Логирование заголовков запроса для отладки
	headers := c.Request.Header
	for name, values := range headers {
		for _, value := range values {
			logrus.Infof("Header: %s=%s", name, value)
		}
	}

	// Получаем пользователя из контекста
	user, exists := c.Get("user")
	if !exists {
		logrus.Warn("User not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	logrus.Infof("User found: %v", user)

	// Проверяем, имеет ли пользователь права администратора
	if !isAdmin(user.(User)) {
		logrus.Warn("User does not have admin privileges")
		c.JSON(http.StatusForbidden, gin.H{"error": "User does not have admin privileges"})
		return
	}

	logrus.Info("User is authorized and has admin privileges")

	utils.HandleRequest(c, func(c *gin.Context) error {
		logrus.Info("Handling request inside AddWalletHandler")
		var wallet WalletStrength
		if err := c.ShouldBindJSON(&wallet); err != nil {
			logrus.Errorf("Invalid request body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return err
		}

		logrus.Infof("Request body: %v", wallet)

		// Проверяем, существует ли уже запись с таким wallet_address
		if _, err := repository.GetVoteStrength(wallet.WalletAddress); err == nil {
			logrus.Errorf("Wallet address already exists: %v", wallet.WalletAddress)
			c.JSON(http.StatusConflict, gin.H{"error": "Wallet address already exists"})
			return err
		}

		logrus.Infof("Wallet address not found, proceeding to add: %v", wallet.WalletAddress)
		if err := repository.AddWalletStrength(wallet.WalletAddress, wallet.VotePower); err != nil {
			logrus.Errorf("Failed to add wallet: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add wallet"})
			return err
		}

		logrus.Info("Wallet added successfully")
		c.JSON(http.StatusCreated, gin.H{"message": "Wallet added successfully"})
		logrus.Info("AddWalletHandler completed successfully")
		return nil
	})
}

// DeleteWalletHandler удаляет адрес кошелька и силу голоса из базы данных
func DeleteWalletHandler(c *gin.Context) {
	logrus.Info("DeleteWalletHandler called")
	user, exists := c.Get("user")
	if !exists {
		logrus.Warn("User not authorized")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	logrus.Infof("User found: %v", user)
	if !isAdmin(user.(User)) {
		logrus.Warn("User does not have admin privileges")
		c.JSON(http.StatusForbidden, gin.H{"error": "User does not have admin privileges"})
		return
	}

	logrus.Info("User is authorized and has admin privileges")
	utils.HandleRequest(c, func(c *gin.Context) error {
		logrus.Info("Handling request inside DeleteWalletHandler")
		walletAddress := c.Param("wallet_address")
		logrus.Infof("Wallet address to delete: %v", walletAddress)
		if err := repository.DeleteWalletStrength(walletAddress); err != nil {
			logrus.Errorf("Failed to delete wallet: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete wallet"})
			return err
		}

		logrus.Info("Wallet deleted successfully")
		c.JSON(http.StatusOK, gin.H{"message": "Wallet deleted successfully"})
		logrus.Info("DeleteWalletHandler completed successfully")
		return nil
	})
}

// GetTableNamesHandler Обработчик маршрута для получения названий таблиц
func GetTableNamesHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	if !isAdmin(user.(User)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "User does not have admin privileges"})
		return
	}

	utils.HandleRequest(c, func(c *gin.Context) error {
		tableNames, err := repository.GetTableNames()
		if err != nil {
			logrus.Errorf("Failed to get table names: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get table names"})
			return err
		}

		c.JSON(http.StatusOK, gin.H{"table_names": tableNames})
		return nil
	})
}

// GetTableElementsHandler обработчик маршрута для получения элементов в таблице по названию
func GetTableElementsHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authorized"})
		return
	}

	if !isAdmin(user.(User)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "User does not have admin privileges"})
		return
	}

	tableName := c.Param("table_name")
	utils.HandleRequest(c, func(c *gin.Context) error {
		elements, err := repository.GetTableElements(tableName)
		if err != nil {
			logrus.Errorf("Failed to get elements from table %s: %v", tableName, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get elements from table"})
			return err
		}

		c.JSON(http.StatusOK, gin.H{"elements": elements})
		return nil
	})
}
