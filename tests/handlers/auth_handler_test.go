package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dao_vote/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Тестирование сервиса авторизации

// mockAuthRequest представляет тестовый запрос для авторизации
var mockAuthRequest = handlers.AuthRequest{
	Login:      "aleksei.ikt@gmail.com",
	Password:   "123qweasd",
	DeviceName: "mobile",
}

// TestUserLoginHandler тестирует обработчик UserLoginHandler
func TestUserLoginHandler(t *testing.T) {
	// Создаем новый Gin роутер и регистрируем обработчик
	router := gin.Default()
	router.POST("/auth/login", handlers.UserLoginHandler)

	// Преобразуем запрос в JSON
	requestBody, _ := json.Marshal(mockAuthRequest)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для получения ответа
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем статус код и ответ
	assert.Equal(t, http.StatusOK, w.Code)

	var response handlers.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Token)

}

// TestUserMeHandler тестирует обработчик UserMeHandler
func TestUserMeHandler(t *testing.T) {
	router := gin.Default()
	router.GET("/auth/me", handlers.UserMeHandler)

	// Создаем новый запрос с  токеном
	req, _ := http.NewRequest("GET", "/auth/me", nil)
	req.Header.Set("Authorization", "Bearer 1825|oyVzunuVE1tuwTmkkOGCfiijz9hT9nJY5fX9O7Xp")
	req.Header.Set("Content-Type", "application/json")

	// Создаем ResponseRecorder для получения ответа
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем статус код и ответ
	assert.Equal(t, http.StatusOK, w.Code)

	var response handlers.UserMeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Data.ID)
}
