package handlers_test

import (
	"bytes"
	"dao_vote/internal/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mockAuthRequest представляет тестовый запрос для авторизации
var mockAuthRequest = handlers.AuthRequest{
	Login:      "aleksei.ikt@gmail.com ",
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
