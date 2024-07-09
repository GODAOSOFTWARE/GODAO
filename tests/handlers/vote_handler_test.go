package handlers_test

import (
	"bytes"
	"dao_vote/internal/handlers"
	"dao_vote/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var mockVoteRequest = url.Values{
	"title":       {"Голосование №1"},
	"subtitle":    {"Суть предложения"},
	"description": {"Подробное описание"},
	"choice":      {"За"},
}

var mockUserVoteRequest = models.UserVote{
	VoteID:    1,
	Voter:     "d01p55v08ld8yc0my72ccpsztv7auyxn2tden6yvw",
	Choice:    "За",
	VotePower: 1000000,
}

func TestCreateVoteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/votes", handlers.CreateVoteHandler)

	// Создаем новый запрос с данными формы
	req, _ := http.NewRequest("POST", "/votes", strings.NewReader(mockVoteRequest.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer 1825|oyVzunuVE1tuwTmkkOGCfiijz9hT9nJY5fX9O7Xp")

	// Создаем ResponseRecorder для получения ответа
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем статус код и ответ
	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Vote
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.ID)
}

func TestGetVoteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/votes/:id", handlers.GetVoteHandler)

	req, _ := http.NewRequest("GET", "/votes/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Vote
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.ID)
}

func TestDeleteVoteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/votes/:id", handlers.DeleteVoteHandler)

	req, _ := http.NewRequest("DELETE", "/votes/1", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestAddUserVoteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/votes/:id/vote", handlers.AddUserVoteHandler)

	requestBody, _ := json.Marshal(mockUserVoteRequest)
	req, _ := http.NewRequest("POST", "/votes/1/vote", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer testtoken")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["transaction_id"])
	assert.NotEmpty(t, response["transaction_hash"])
}
