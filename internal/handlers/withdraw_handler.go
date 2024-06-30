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
	"time"
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

	// Логирование токена и запроса на снятие средств
	logrus.Infof("Authorization token: %s", token)
	logrus.Infof("Withdraw request data: %s", jsonData)

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

	logrus.Infof("Sending request to external API for withdrawal")
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

	// Логирование ответа на запрос снятия средств
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
			TransactionID int `json:"transaction_id"`
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

	// Функция ожидания хэша транзакции
	getTransactionHash := func(transactionID int) (string, error) {
		hashReqURL := fmt.Sprintf("https://backend.ddapps.io/api/v1/transactions/%d", transactionID)
		for i := 0; i < 10; i++ {
			hashReq, err := http.NewRequest("GET", hashReqURL, nil)
			if err != nil {
				logrus.Errorf("Failed to create new request for transaction hash: %v", err)
				return "", err
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
				return "", err
			}
			defer hashResp.Body.Close()

			// Чтение и обработка ответа от внешнего API
			hashBody, err := ioutil.ReadAll(hashResp.Body)
			if err != nil {
				logrus.Errorf("Failed to read response body for transaction hash: %v", err)
				return "", err
			}

			// Логирование ответа на запрос хэша транзакции
			logrus.Infof("Response from external API for transaction hash: %s", string(hashBody))
			if hashResp.StatusCode != http.StatusOK {
				logrus.Errorf("Error from external API for transaction hash: %s", string(hashBody))
				return "", fmt.Errorf("error from external API: %s", string(hashBody))
			}

			var hashResponse struct {
				Data struct {
					Hash string `json:"hash"`
				} `json:"data"`
			}
			if err := json.Unmarshal(hashBody, &hashResponse); err != nil {
				logrus.Errorf("Failed to unmarshal response JSON for transaction hash: %v", err)
				return "", err
			}

			transactionHash := hashResponse.Data.Hash
			logrus.Infof("Transaction hash: %s", transactionHash)

			if transactionHash != "" {
				return transactionHash, nil
			}

			// Задержка перед повторной попыткой
			time.Sleep(1 * time.Second)
		}

		return "", fmt.Errorf("transaction hash not found after multiple attempts")
	}

	transactionHash, err := getTransactionHash(transactionID)
	if err != nil {
		logrus.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful, but failed to retrieve transaction hash"})
		return
	}

	// Успешный ответ с хэшем транзакции
	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful", "transaction_id": transactionID, "transaction_hash": transactionHash})
}
