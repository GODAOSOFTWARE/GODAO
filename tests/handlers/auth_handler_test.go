package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dao_vote/back-end/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mockAuthRequest представляет тестовый запрос для авторизации
var mockAuthRequest = handlers.AuthRequest{
	Login:      "email",                    // Тестовый логин
	Password:   "password (min 6 letters)", // Тестовый пароль
	DeviceName: "your device",              // Тестовое имя устройства
}

// TestUserLoginHandler тестирует обработчик UserLoginHandler
func TestUserLoginHandler(t *testing.T) {
	router := gin.Default()                               // Создаем новый Gin роутер
	router.POST("/auth/login", handlers.UserLoginHandler) // Регистрируем обработчик для маршрута POST /auth/login

	requestBody, _ := json.Marshal(mockAuthRequest)                                // Преобразуем тестовый запрос в JSON
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(requestBody)) // Создаем новый HTTP запрос с JSON телом
	req.Header.Set("Content-Type", "application/json")                             // Устанавливаем заголовок Content-Type

	w := httptest.NewRecorder() // Создаем ResponseRecorder для записи ответа
	router.ServeHTTP(w, req)    // Передаем запрос в роутер

	assert.Equal(t, http.StatusOK, w.Code) // Проверяем, что статус код ответа 200 OK

	var response handlers.AuthResponse               // Объявляем переменную для хранения JSON ответа
	err := json.Unmarshal(w.Body.Bytes(), &response) // Распаковываем JSON ответ в переменную
	assert.NoError(t, err)                           // Проверяем, что при распаковке не возникло ошибок
	assert.NotEmpty(t, response.Token)               // Проверяем, что в ответе присутствует поле token
}

// TestUserMeHandler тестирует обработчик UserMeHandler
func TestUserMeHandler(t *testing.T) {
	router := gin.Default()                        // Создаем новый Gin роутер
	router.GET("/auth/me", handlers.UserMeHandler) // Регистрируем обработчик для маршрута GET /auth/me

	req, _ := http.NewRequest("GET", "/auth/me", nil)  // Создаем новый HTTP GET запрос
	req.Header.Set("Authorization", "Bearer ")         // Устанавливаем заголовок Authorization
	req.Header.Set("Content-Type", "application/json") // Устанавливаем заголовок Content-Type

	w := httptest.NewRecorder() // Создаем ResponseRecorder для записи ответа
	router.ServeHTTP(w, req)    // Передаем запрос в роутер

	assert.Equal(t, http.StatusOK, w.Code) // Проверяем, что статус код ответа 200 OK

	var response handlers.UserMeResponse             // Объявляем переменную для хранения JSON ответа
	err := json.Unmarshal(w.Body.Bytes(), &response) // Распаковываем JSON ответ в переменную
	assert.NoError(t, err)                           // Проверяем, что при распаковке не возникло ошибок
	assert.NotEmpty(t, response.Data.ID)             // Проверяем, что в ответе присутствует поле ID
}
