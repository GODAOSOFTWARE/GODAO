// tests/handlers/vote_handler_test.go
// update
package handlers_test

import (
	"bytes"
	"dao_vote/internal/handlers"
	"dao_vote/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var mockVoteRequest = map[string]string{
	"title":       "Test Vote",
	"subtitle":    "Test Subtitle",
	"description": "Test Description",
	"choice":      "За",
}

var mockUserVoteRequest = models.UserVote{
	VoteID:    1,
	Voter:     "test_voter",
	Choice:    "За",
	VotePower: 100,
}

func TestCreateVoteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/votes", handlers.CreateVoteHandler)

	requestBody, _ := json.Marshal(mockVoteRequest)
	req, _ := http.NewRequest("POST", "/votes", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer testtoken")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

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
