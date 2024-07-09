package handlers_test

import (
	"bytes"
	"dao_vote/back-end/handlers"
	"dao_vote/back-end/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mockWithdrawRequest представляет тестовый запрос для вывода средств
var mockWithdrawRequest = models.WithdrawRequest{
	Amount:  1.0,                                         // Устанавливаем тестовую сумму вывода
	Address: "d01juva4qeqjyavwaf4s2vfzpg2y8vj6gl9dtne45", // Устанавливаем тестовый адрес кошелька
}

// TestWithdrawHandler тестирует обработчик WithdrawHandler
func TestWithdrawHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)                                 // Устанавливаем режим тестирования Gin
	router := gin.Default()                                   // Создаем новый роутер Gin
	router.POST("/api/v1/withdraw", handlers.WithdrawHandler) // Регистрируем обработчик для маршрута POST /api/v1/withdraw

	requestBody, _ := json.Marshal(mockWithdrawRequest)                                 // Преобразуем тестовый запрос в JSON
	req, _ := http.NewRequest("POST", "/api/v1/withdraw", bytes.NewBuffer(requestBody)) // Создаем новый HTTP запрос с JSON телом
	req.Header.Set("Content-Type", "application/json")                                  // Устанавливаем заголовок Content-Type
	req.Header.Set("Authorization", "Bearer ")                                          // Устанавливаем заголовок Authorization

	w := httptest.NewRecorder() // Создаем ResponseRecorder для записи ответа
	router.ServeHTTP(w, req)    // Передаем запрос в роутер

	assert.Equal(t, http.StatusOK, w.Code) // Проверяем, что статус код ответа 200 OK

	var response map[string]interface{}              // Объявляем переменную для хранения JSON ответа
	err := json.Unmarshal(w.Body.Bytes(), &response) // Распаковываем JSON ответ в переменную
	assert.NoError(t, err)                           // Проверяем, что при распаковке не возникло ошибок
	assert.NotEmpty(t, response["transaction_id"])   // Проверяем, что в ответе присутствует поле transaction_id
}
