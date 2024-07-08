// tests/handlers/withdraw_handler_test.go
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

var mockWithdrawRequest = models.WithdrawRequest{
	Amount:  1.0,
	Address: "d01juva4qeqjyavwaf4s2vfzpg2y8vj6gl9dtne45",
}

func TestWithdrawHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/api/v1/withdraw", handlers.WithdrawHandler)

	requestBody, _ := json.Marshal(mockWithdrawRequest)
	req, _ := http.NewRequest("POST", "/api/v1/withdraw", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer 1825|oyVzunuVE1tuwTmkkOGCfiijz9hT9nJY5fX9O7Xp")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["transaction_id"])
}
