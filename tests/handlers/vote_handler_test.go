package handlers_test

import (
	"bytes"
	"dao_vote/back-end/handlers"
	"dao_vote/back-end/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mockVoteRequest представляет тестовый запрос для создания голосования
var mockVoteRequest = url.Values{
	"title":       {"Голосование №1"},     // Тестовый заголовок
	"subtitle":    {"Суть предложения"},   // Тестовый подзаголовок
	"description": {"Подробное описание"}, // Тестовое описание
	"choice":      {"За"},                 // Тестовый выбор
}

// mockUserVoteRequest представляет тестовый запрос для голосования пользователя
var mockUserVoteRequest = models.UserVote{
	VoteID:    1,                                           // Тестовый ID голосования
	Voter:     "d01p55v08ld8yc0my72ccpsztv7auyxn2tden6yvw", // Тестовый адрес кошелька
	Choice:    "За",                                        // Тестовый выбор
	VotePower: 1000000,                                     // Тестовая сила голоса
}

// TestCreateVoteHandler тестирует обработчик CreateVoteHandler
func TestCreateVoteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)                         // Устанавливаем режим тестирования Gin
	router := gin.Default()                           // Создаем новый роутер Gin
	router.POST("/votes", handlers.CreateVoteHandler) // Регистрируем обработчик для маршрута POST /votes

	req, _ := http.NewRequest("POST", "/votes", strings.NewReader(mockVoteRequest.Encode())) // Создаем новый HTTP запрос с данными формы
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")                      // Устанавливаем заголовок Content-Type
	req.Header.Set("Authorization", "Bearer ")                                               // Устанавливаем заголовок Authorization

	w := httptest.NewRecorder() // Создаем ResponseRecorder для записи ответа
	router.ServeHTTP(w, req)    // Передаем запрос в роутер

	assert.Equal(t, http.StatusCreated, w.Code) // Проверяем, что статус код ответа 201 Created

	var response models.VoteInfo                     // Объявляем переменную для хранения JSON ответа
	err := json.Unmarshal(w.Body.Bytes(), &response) // Распаковываем JSON ответ в переменную
	assert.NoError(t, err)                           // Проверяем, что при распаковке не возникло ошибок
	assert.NotEmpty(t, response.ID)                  // Проверяем, что в ответе присутствует поле ID
}

// TestGetVoteHandler тестирует обработчик GetVoteHandler
func TestGetVoteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)                         // Устанавливаем режим тестирования Gin
	router := gin.Default()                           // Создаем новый роутер Gin
	router.GET("/votes/:id", handlers.GetVoteHandler) // Регистрируем обработчик для маршрута GET /votes/:id

	req, _ := http.NewRequest("GET", "/votes/1", nil) // Создаем новый HTTP GET запрос

	w := httptest.NewRecorder() // Создаем ResponseRecorder для записи ответа
	router.ServeHTTP(w, req)    // Передаем запрос в роутер

	assert.Equal(t, http.StatusOK, w.Code) // Проверяем, что статус код ответа 200 OK

	var response models.VoteInfo                     // Объявляем переменную для хранения JSON ответа
	err := json.Unmarshal(w.Body.Bytes(), &response) // Распаковываем JSON ответ в переменную
	assert.NoError(t, err)                           // Проверяем, что при распаковке не возникло ошибок
	assert.NotEmpty(t, response.ID)                  // Проверяем, что в ответе присутствует поле ID
}

// TestDeleteVoteHandler тестирует обработчик DeleteVoteHandler
func TestDeleteVoteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)                               // Устанавливаем режим тестирования Gin
	router := gin.Default()                                 // Создаем новый роутер Gin
	router.DELETE("/votes/:id", handlers.DeleteVoteHandler) // Регистрируем обработчик для маршрута DELETE /votes/:id

	req, _ := http.NewRequest("DELETE", "/votes/1", nil) // Создаем новый HTTP DELETE запрос

	w := httptest.NewRecorder() // Создаем ResponseRecorder для записи ответа
	router.ServeHTTP(w, req)    // Передаем запрос в роутер

	assert.Equal(t, http.StatusNoContent, w.Code) // Проверяем, что статус код ответа 204 No Content
}

// TestAddUserVoteHandler тестирует обработчик AddUserVoteHandler
func TestAddUserVoteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)                                   // Устанавливаем режим тестирования Gin
	router := gin.Default()                                     // Создаем новый роутер Gin
	router.POST("/votes/:id/vote", handlers.AddUserVoteHandler) // Регистрируем обработчик для маршрута POST /votes/:id/vote

	requestBody, _ := json.Marshal(mockUserVoteRequest)                              // Преобразуем тестовый запрос в JSON
	req, _ := http.NewRequest("POST", "/votes/1/vote", bytes.NewBuffer(requestBody)) // Создаем новый HTTP POST запрос с JSON телом
	req.Header.Set("Content-Type", "application/json")                               // Устанавливаем заголовок Content-Type
	req.Header.Set("Authorization", "Bearer testtoken")                              // Устанавливаем заголовок Authorization

	w := httptest.NewRecorder() // Создаем ResponseRecorder для записи ответа
	router.ServeHTTP(w, req)    // Передаем запрос в роутер

	assert.Equal(t, http.StatusOK, w.Code) // Проверяем, что статус код ответа 200 OK

	var response map[string]interface{}              // Объявляем переменную для хранения JSON ответа
	err := json.Unmarshal(w.Body.Bytes(), &response) // Распаковываем JSON ответ в переменную
	assert.NoError(t, err)                           // Проверяем, что при распаковке не возникло ошибок
	assert.NotEmpty(t, response["transaction_id"])   // Проверяем, что в ответе присутствует поле transaction_id
	assert.NotEmpty(t, response["transaction_hash"]) // Проверяем, что в ответе присутствует поле transaction_hash
}
