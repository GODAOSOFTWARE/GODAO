package services

import (
	"dao_vote/internal/models"
	"dao_vote/internal/repository"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	requiredMajority = 51  // Требуемое большинство для принятия решения
	percentFactor    = 100 // Фактор для расчета процентов
)

// FetchDAOTeamVoteResults - функция для получения результатов голосования по адресу кошелька DAO
func FetchDAOTeamVoteResults(walletAddress string, offset int) (models.DAOTeamApiResponse, error) {
	// Получаем значение лимита из числа членов DAO
	limit := len(repository.GetVoteMap())

	// Формируем URL для запроса с динамическими параметрами limit и offset
	apiURL := fmt.Sprintf("https://mainnet-explorer-api.decimalchain.com/api/address/%s/txs?limit=%d&offset=%d", walletAddress, limit, offset)

	// Выполняем GET-запрос к API для получения транзакций голосования
	resp, err := http.Get(apiURL)
	if err != nil {
		return models.DAOTeamApiResponse{}, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем, был ли запрос успешным
	if resp.StatusCode != http.StatusOK {
		return models.DAOTeamApiResponse{}, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	// Читаем тело ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.DAOTeamApiResponse{}, fmt.Errorf("error reading response body: %v", err)
	}

	// Парсим тело ответа в структуру apiResponse
	var apiResponse models.DAOTeamApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return models.DAOTeamApiResponse{}, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	// Обрабатываем каждую транзакцию, добавляя силу голоса
	for i, result := range apiResponse.Result.Txs {
		votePower, err := repository.GetVoteStrength(result.From)
		if err != nil {
			// Логируем ошибку, но продолжаем обработку
			fmt.Printf("Error getting vote strength for %s: %v\n", result.From, err)
		}
		apiResponse.Result.Txs[i].VotePower = votePower
		apiResponse.Result.Txs[i].Hash = result.Hash // Добавляем хэш транзакции
	}

	return apiResponse, nil
}

// PrepareDAOTeamVoteResults - функция для подготовки результатов голосования команды DAO
func PrepareDAOTeamVoteResults(apiResponse models.DAOTeamApiResponse) models.DAOTeamVoteResultsResponse {
	votesFor := make(map[string]bool)
	votesAgainst := make(map[string]bool)

	// Обрабатывает каждую транзакцию, определяя голос "За" или "Против"
	for _, result := range apiResponse.Result.Txs {
		message := strings.TrimSpace(strings.ToLower(result.Message))
		message = strings.Trim(message, `"`)

		switch message {
		case "да", "дa", "д", "за", "зa", "z":
			votesFor[result.From] = true
		case "нет", "против", "н":
			votesAgainst[result.From] = true
		}
	}

	// Расчитывает силу голосов "За" и "Против"
	strengthFor := calculateStrength(votesFor)
	strengthAgainst := calculateStrength(votesAgainst)
	totalVoices := repository.GetTotalVoices()
	percentFor := calculatePercentage(strengthFor, totalVoices)
	percentAgainst := calculatePercentage(strengthAgainst, totalVoices)

	// Определяет статус голосования и резолюцию
	status := "Активно"
	resolution := "Решение не принято"
	if percentFor >= requiredMajority {
		status = "Завершено"
		resolution = "Принять изменения"
	} else if percentAgainst >= requiredMajority {
		status = "Завершено"
		resolution = "Отклонить изменения"
	}

	voteMap := repository.GetVoteMap()

	// Возвращает подготовленные результаты голосования
	return models.DAOTeamVoteResultsResponse{
		DAOMembers:   len(voteMap),
		VotedMembers: len(apiResponse.Result.Txs),
		Turnout:      formatPercentage(float64(len(apiResponse.Result.Txs)) / float64(len(voteMap)) * percentFactor),
		VotesFor:     fmt.Sprintf("%d/%d (%.2f%%)", strengthFor, totalVoices, percentFor),
		VotesAgainst: fmt.Sprintf("%d (%.2f%%)", strengthAgainst, percentAgainst),
		VotingStatus: status,
		Resolution:   resolution,
		Transactions: apiResponse.Result.Txs,
	}
}

// calculateStrength вычисляет общую силу голосов из карты голосов
func calculateStrength(votes map[string]bool) int {
	strength := 0
	for voter := range votes {
		if voteStrength, err := repository.GetVoteStrength(voter); err == nil {
			strength += voteStrength
		} else {
			// Логирует ошибку
			fmt.Printf("Error getting vote strength for %s: %v\n", voter, err)
		}
	}
	return strength
}

// calculatePercentage вычисляет процент от общего числа голосов
func calculatePercentage(value, total int) float64 {
	return float64(value) / float64(total) * percentFactor
}

// formatPercentage форматирует значение процента
func formatPercentage(value float64) string {
	return fmt.Sprintf("%.2f%%", value)
}
