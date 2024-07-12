// Package handlers Обработчик командных голосований
package handlers

import (
	"dao_vote/back-end/services"
	"dao_vote/back-end/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetVotingResultsByWallet обрабатывает GET /dao-team-vote-results запрос
func GetVotingResultsByWallet(c *gin.Context) {
	utils.HandleRequest(c, func(c *gin.Context) error {
		// Получаем параметр wallet_address из запроса
		walletAddress := c.Query("wallet_address")
		if walletAddress == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "wallet_address is required"})
			return nil
		}

		// Получаем параметр offset из запроса
		offsetStr := c.Query("offset")
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			offset = 0 // Устанавливаем значение по умолчанию, если параметр отсутствует или некорректен
		}

		// Вызов сервиса для получения результатов голосования команды DAO
		apiResponse, err := services.FetchDAOTeamVoteResults(walletAddress, offset)
		if err != nil {
			// Возвращает ошибку, если не удалось получить результаты голосования
			return err
		}

		// Подготовка результатов голосования для отправки в ответе
		voteResults := services.PrepareDAOTeamVoteResults(apiResponse)
		// Возвращает успешный ответ с результатами голосования
		c.JSON(http.StatusOK, voteResults)
		return nil
	})
}
