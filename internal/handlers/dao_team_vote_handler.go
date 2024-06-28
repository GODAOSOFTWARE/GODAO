package handlers

import (
	"dao_vote/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetDAOTeamVoteResults обрабатывает GET /dao-team-vote-results запрос
func GetDAOTeamVoteResults(c *gin.Context) {
	// Получаем параметр wallet_address из запроса
	walletAddress := c.Query("wallet_address")
	if walletAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wallet_address is required"})
		return
	}

	// Вызов сервиса для получения результатов голосования команды DAO
	apiResponse, err := services.FetchDAOTeamVoteResults(walletAddress)
	if err != nil {
		// Возвращает ошибку, если не удалось получить результаты голосования
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Подготовка результатов голосования для отправки в ответе
	voteResults := services.PrepareDAOTeamVoteResults(apiResponse)
	// Возвращает успешный ответ с результатами голосования
	c.JSON(http.StatusOK, voteResults)
}
