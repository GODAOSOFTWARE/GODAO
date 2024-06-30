package handlers

import (
	"bytes"
	"dao_vote/internal/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// WithdrawHandler обрабатывает запрос на снятие средств
func WithdrawHandler(c *gin.Context) {
	var withdrawReq models.WithdrawRequest
	if err := c.ShouldBindJSON(&withdrawReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Подготовка запроса к внешнему API
	jsonData, err := json.Marshal(withdrawReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Получение токена из заголовка
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// Отправка запроса к внешнему API
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://backend.ddapps.io/api/v1/withdraw", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Чтение и обработка ответа от внешнего API
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Error from external API: %s", string(body))
		c.JSON(resp.StatusCode, gin.H{"error": string(body)})
		return
	}

	// Возвращение ответа клиенту
	c.JSON(http.StatusOK, gin.H{"message": "Withdrawal successful"})
}
