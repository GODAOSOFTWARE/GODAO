package handlers

import (
	"bytes"
	"dao_vote/internal/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// WithdrawHandler обрабатывает запрос на снятие средств
func WithdrawHandler(c *gin.Context) {
	var withdrawReq models.WithdrawRequest
	if err := c.ShouldBindJSON(&withdrawReq); err != nil {
		logrus.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Подготовка запроса к внешнему API
	jsonData, err := json.Marshal(withdrawReq)
	if err != nil {
		logrus.Errorf("Failed to marshal JSON: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Получение токена из заголовка
	token := c.GetHeader("Authorization")
	if token == "" {
		logrus.Error("Authorization token is missing")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// Отправка запроса на снятие средств к внешнему API
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://backend.ddapps.io/api/v1/withdraw", bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.Errorf("Failed to create new request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	logrus.Infof("Sending request to external API for withdrawal: %v", withdrawReq)
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Failed to send request to external API: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Чтение и обработка ответа от внешнего API
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.Infof("Response from external API: %s", string(body))
	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Error from external API: %s", string(body))
		c.JSON(resp.StatusCode, gin.H{"error": string(body)})
		return
	}

	var withdrawResponse struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		Data    struct {
			TransactionID float64 `json:"transaction_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &withdrawResponse); err != nil {
		logrus.Errorf("Failed to unmarshal response JSON: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.Infof("Withdraw response JSON: %v", withdrawResponse)
	transactionID := withdrawResponse.Data.TransactionID
	if transactionID == 0 {
		logrus.Error("Invalid transaction ID in response")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid transaction ID in response"})
		return
	}

	// Запрос хэша транзакции
	hashReqURL := fmt.Sprintf("https://backend.ddapps.io/api/v1/transaction/%d/hash", int(transactionID))
	hashReq, err := http.NewRequest("GET", hashReqURL, nil)
	if err != nil {
		logrus.Errorf("Failed to create new request for transaction hash: %v", err)
		// Возвращаем успешный ответ, несмотря на ошибку запроса хэша
		c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful, but failed to retrieve transaction hash"})
		return
	}
	hashReq.Header.Set("Authorization", token)
	hashReq.Header.Set("Accept", "application/json")

	logrus.Infof("Sending request to external API for transaction hash: %v", hashReqURL)
	hashResp, err := client.Do(hashReq)
	if err != nil {
		logrus.Errorf("Failed to send request to external API for transaction hash: %v", err)
		// Возвращаем успешный ответ, несмотря на ошибку запроса хэша
		c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful, but failed to retrieve transaction hash"})
		return
	}
	defer hashResp.Body.Close()

	// Чтение и обработка ответа от внешнего API
	hashBody, err := ioutil.ReadAll(hashResp.Body)
	if err != nil {
		logrus.Errorf("Failed to read response body for transaction hash: %v", err)
		// Возвращаем успешный ответ, несмотря на ошибку запроса хэша
		c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful, but failed to retrieve transaction hash"})
		return
	}

	logrus.Infof("Response from external API for transaction hash: %s", string(hashBody))
	if hashResp.StatusCode != http.StatusOK {
		logrus.Errorf("Error from external API for transaction hash: %s", string(hashBody))
		// Возвращаем успешный ответ, несмотря на ошибку запроса хэша
		c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful, but failed to retrieve transaction hash"})
		return
	}

	var hashResponse struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		Data    struct {
			Hash string `json:"hash"`
		} `json:"data"`
	}
	if err := json.Unmarshal(hashBody, &hashResponse); err != nil {
		logrus.Errorf("Failed to unmarshal response JSON for transaction hash: %v", err)
		// Возвращаем успешный ответ, несмотря на ошибку запроса хэша
		c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful, but failed to retrieve transaction hash"})
		return
	}

	logrus.Infof("Transaction hash response JSON: %v", hashResponse)
	transactionHash := hashResponse.Data.Hash
	if transactionHash == "" {
		logrus.Error("Invalid transaction hash in response")
		// Возвращаем успешный ответ, несмотря на ошибку запроса хэша
		c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful, but failed to retrieve transaction hash"})
		return
	}

	// Успешный ответ с хэшем транзакции
	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful", "transaction_id": transactionID, "transaction_hash": transactionHash})
}
