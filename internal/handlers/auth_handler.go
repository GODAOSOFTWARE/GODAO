package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

// AuthRequest представляет структуру запроса для авторизации
type AuthRequest struct {
	Login      string `json:"login"`
	Password   string `json:"password"`
	DeviceName string `json:"device_name"`
}

// AuthResponse представляет структуру ответа для авторизации
type AuthResponse struct {
	Token string `json:"token"`
}

// UserLoginHandler обрабатывает запрос авторизации пользователя
func UserLoginHandler(c *gin.Context) {
	var authReq AuthRequest
	if err := c.ShouldBindJSON(&authReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Подготовка запроса к внешнему API
	jsonData, err := json.Marshal(authReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Отправка запроса к внешнему API
	resp, err := http.Post("https://backend.ddapps.io/api/v1/auth/user_login", "application/json", bytes.NewBuffer(jsonData))
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

	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращение ответа клиенту
	c.JSON(http.StatusOK, authResp)
}
