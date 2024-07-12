// Package handlers Обработчик снятия токенов с баланса
package handlers

import (
	"bytes"
	"dao_vote/back-end/models"
	"dao_vote/back-end/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

// WithdrawHandler обрабатывает запрос на снятие средств
func WithdrawHandler(c *gin.Context) {
	utils.HandleRequest(c, func(c *gin.Context) error {
		var withdrawReq models.WithdrawRequest
		if err := c.ShouldBindJSON(&withdrawReq); err != nil {
			logrus.Errorf("Failed to bind JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return nil
		}

		// Подготовка запроса к внешнему API
		jsonData, err := json.Marshal(withdrawReq)
		if err != nil {
			logrus.Errorf("Failed to marshal JSON: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return nil
		}

		// Получение токена из заголовка
		token := c.GetHeader("Authorization")
		if token == "" {
			logrus.Error("Authorization token is missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
			return nil
		}

		// Логирование токена и запроса на снятие средств
		logrus.Infof("Authorization token: %s", token)
		logrus.Infof("Withdraw request data: %s", jsonData)

		// Отправка запроса на снятие средств к внешнему API
		client := &http.Client{}
		req, err := http.NewRequest("POST", "https://backend.ddapps.io/api/v1/withdraw", bytes.NewBuffer(jsonData))
		if err != nil {
			logrus.Errorf("Failed to create new request: %v", err)
			return err
		}

		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		logrus.Infof("Sending request to external API for withdrawal")
		resp, err := client.Do(req)
		if err != nil {
			logrus.Errorf("Failed to send request to external API: %v", err)
			return err
		}
		defer resp.Body.Close()

		// Чтение и обработка ответа от внешнего API
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("Failed to read response body: %v", err)
			return err
		}

		// Логирование ответа на запрос снятия средств
		logrus.Infof("Response from external API: %s", string(body))
		if resp.StatusCode != http.StatusOK {
			logrus.Errorf("Error from external API: %s", string(body))
			c.JSON(resp.StatusCode, gin.H{"error": string(body)})
			return nil
		}

		var withdrawResponse struct {
			Type    string `json:"type"`
			Message string `json:"message"`
			Data    struct {
				TransactionID int `json:"transaction_id"`
			} `json:"data"`
		}
		if err := json.Unmarshal(body, &withdrawResponse); err != nil {
			logrus.Errorf("Failed to unmarshal response JSON: %v", err)
			return err
		}

		logrus.Infof("Withdraw response JSON: %v", withdrawResponse)
		transactionID := withdrawResponse.Data.TransactionID
		if transactionID == 0 {
			logrus.Error("Invalid transaction ID in response")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid transaction ID in response"})
			return nil
		}

		// Задержка перед запросом хэша транзакции
		time.Sleep(5 * time.Second) // Задержка 5 секунд

		// Запрос хэша транзакции
		hashReqURL := fmt.Sprintf("https://backend.ddapps.io/api/v1/transactions/%d", transactionID)
		hashReq, err := http.NewRequest("GET", hashReqURL, nil)
		if err != nil {
			logrus.Errorf("Failed to create new request for transaction hash: %v", err)
			// Возвращаем успешный ответ, несмотря на ошибку запроса хэша
			c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful, but failed to retrieve transaction hash"})
			return nil
		}
		hashReq.Header.Set("Authorization", token)
		hashReq.Header.Set("Accept", "application/json")

		// Логирование запроса на получение хэша транзакции
		logrus.Infof("Sending request to external API for transaction hash: %v", hashReqURL)
		logrus.Infof("Request URL: %v", hashReqURL)
		logrus.Infof("Request method: %v", hashReq.Method)
		logrus.Infof("Request headers: %v", hashReq.Header)
		logrus.Infof("Request ID: %d", transactionID)

		hashResp, err := client.Do(hashReq)
		if err != nil {
			logrus.Errorf("Failed to send request to external API for transaction hash: %v", err)
			// Возвращаем успешный ответ, несмотря на ошибку запроса хэша
			c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful, but failed to retrieve transaction hash"})
			return nil
		}
		defer hashResp.Body.Close()

		// Чтение и обработка ответа от внешнего API
		hashBody, err := ioutil.ReadAll(hashResp.Body)
		if err != nil {
			logrus.Errorf("Failed to read response body for transaction hash: %v", err)
			// Возвращаем успешный ответ, несмотря на ошибку запроса хэша
			c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful, but failed to retrieve transaction hash"})
			return nil
		}

		// Логирование ответа на запрос хэша транзакции
		logrus.Infof("Response from external API for transaction hash: %s", string(hashBody))
		if hashResp.StatusCode != http.StatusOK {
			logrus.Errorf("Error from external API for transaction hash: %s", string(hashBody))
			// Возвращаем успешный ответ, несмотря на ошибку запроса хэша
			c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful, but failed to retrieve transaction hash"})
			return nil
		}

		var hashResponse struct {
			Data struct {
				ID             int    `json:"id"`
				UserID         int    `json:"user_id"`
				WalletID       int    `json:"wallet_id"`
				Amount         string `json:"amount"`
				Payload        string `json:"payload"`
				DecimalPayload string `json:"decimal_payload"`
				Coin           string `json:"coin"`
				Hash           string `json:"hash"`
				Result         string `json:"result"`
				Complete       bool   `json:"complete"`
				Success        bool   `json:"success"`
				CreatedAt      string `json:"created_at"`
				UpdatedAt      string `json:"updated_at"`
				Type           int    `json:"type"`
				Calculated     int    `json:"calculated"`
			} `json:"data"`
		}
		if err := json.Unmarshal(hashBody, &hashResponse); err != nil {
			logrus.Errorf("Failed to unmarshal response JSON for transaction hash: %v", err)
			// Возвращаем успешный ответ, несмотря на ошибку запроса хэша
			c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful, but failed to retrieve transaction hash"})
			return nil
		}

		// Логирование значений полей структуры hashResponse
		logrus.Infof("Transaction hash response data: %+v", hashResponse.Data)

		// Явное логирование значения хэша
		transactionHash := hashResponse.Data.Hash
		logrus.Infof("Transaction hash: %s", transactionHash)

		if transactionHash == "" {
			logrus.Error("Invalid transaction hash in response")
			// Возвращаем успешный ответ, несмотря на ошибку запроса хэша
			c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful, but failed to retrieve transaction hash"})
			return nil
		}

		// Успешный ответ с хэшем транзакции
		c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful", "transaction_id": transactionID, "transaction_hash": transactionHash})
		return nil
	})
}
