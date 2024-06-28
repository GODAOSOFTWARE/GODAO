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
	limit := len(repository.GetVoteMap())
	apiURL := fmt.Sprintf("https://mainnet-explorer-api.decimalchain.com/api/address/%s/txs?limit=%d&offset=%d", walletAddress, limit, offset)

	resp, err := http.Get(apiURL)
	if err != nil {
		return models.DAOTeamApiResponse{}, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.DAOTeamApiResponse{}, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.DAOTeamApiResponse{}, fmt.Errorf("error reading response body: %v", err)
	}

	var apiResponse models.DAOTeamApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return models.DAOTeamApiResponse{}, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	for i, result := range apiResponse.Result.Txs {
		votePower, err := repository.GetVoteStrength(result.From)
		if err != nil {
			fmt.Printf("Error getting vote strength for %s: %v\n", result.From, err)
		}
		apiResponse.Result.Txs[i].VotePower = votePower
		apiResponse.Result.Txs[i].Hash = result.Hash
	}

	return apiResponse, nil
}

// PrepareDAOTeamVoteResults - функция для подготовки результатов голосования команды DAO
func PrepareDAOTeamVoteResults(apiResponse models.DAOTeamApiResponse) models.DAOTeamVoteResultsResponse {
	votesFor := make(map[string]bool)                // Карта голосов "За"
	votesAgainst := make(map[string]bool)            // Карта голосов "Против"
	uniqueVoters := make(map[string]bool)            // Карта уникальных проголосовавших с ненулевой силой голоса
	rejectedTxs := []models.Transaction{}            // Список отклоненных транзакций
	nullVotePowerTxs := []models.Transaction{}       // Список транзакций с нулевой силой голоса
	totalTransactions := len(apiResponse.Result.Txs) // Общее количество транзакций

	for _, result := range apiResponse.Result.Txs {
		// Логируем информацию о текущей транзакции
		fmt.Printf("Processing transaction from %s with vote power %d and message '%s'\n", result.From, result.VotePower, result.Message)

		if uniqueVoters[result.From] {
			// Если пользователь уже голосовал, добавляем транзакцию в список отклоненных
			fmt.Printf("Duplicate transaction detected from %s\n", result.From)
			rejectedTxs = append(rejectedTxs, result)
			continue
		}

		if result.VotePower == 0 {
			// Добавляем транзакцию с нулевой силой голоса в соответствующий список
			fmt.Printf("Null vote power transaction from %s added to nullVotePowerTxs\n", result.From)
			nullVotePowerTxs = append(nullVotePowerTxs, result)
			continue
		}

		// Если транзакция валидная, помечаем пользователя как проголосовавшего
		uniqueVoters[result.From] = true

		message := strings.TrimSpace(strings.ToLower(result.Message))
		message = strings.Trim(message, `"`)

		switch message {
		case "да", "дa", "д", "за", "зa", "z":
			fmt.Printf("Vote FOR detected from %s\n", result.From)
			votesFor[result.From] = true
		case "нет", "против", "н":
			fmt.Printf("Vote AGAINST detected from %s\n", result.From)
			votesAgainst[result.From] = true
		}
	}

	strengthFor := calculateStrength(votesFor)
	strengthAgainst := calculateStrength(votesAgainst)
	totalVoices := repository.GetTotalVoices()
	percentFor := calculatePercentage(strengthFor, totalVoices)
	percentAgainst := calculatePercentage(strengthAgainst, totalVoices)

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

	// Логируем итоговые результаты
	fmt.Printf("Voting completed with %d DAO members, %d voted members\n", len(voteMap), len(votesFor)+len(votesAgainst))
	fmt.Printf("Votes for: %d/%d (%.2f%%), Votes against: %d (%.2f%%)\n", strengthFor, totalVoices, percentFor, strengthAgainst, percentAgainst)
	fmt.Printf("Voting status: %s, Resolution: %s\n", status, resolution)

	return models.DAOTeamVoteResultsResponse{
		DAOMembers:        len(voteMap),
		VotedMembers:      len(votesFor) + len(votesAgainst),
		Turnout:           formatPercentage(float64(len(votesFor)+len(votesAgainst)) / float64(len(voteMap)) * percentFactor),
		VotesFor:          fmt.Sprintf("%d/%d (%.2f%%)", strengthFor, totalVoices, percentFor),
		VotesAgainst:      fmt.Sprintf("%d (%.2f%%)", strengthAgainst, percentAgainst),
		VotingStatus:      status,
		Resolution:        resolution,
		Transactions:      apiResponse.Result.Txs,
		TotalTransactions: totalTransactions,
		RejectedTxs:       rejectedTxs,
		NullVotePowerTxs:  nullVotePowerTxs,
	}
}

// calculateStrength - функция для вычисления общей силы голосов из карты голосов
func calculateStrength(votes map[string]bool) int {
	strength := 0
	for voter := range votes {
		if voteStrength, err := repository.GetVoteStrength(voter); err == nil {
			strength += voteStrength
		} else {
			fmt.Printf("Error getting vote strength for %s: %v\n", voter, err)
		}
	}
	return strength
}

// calculatePercentage - функция для вычисления процента от общего числа голосов
func calculatePercentage(value, total int) float64 {
	return float64(value) / float64(total) * percentFactor
}

// formatPercentage - функция для форматирования значения процента
func formatPercentage(value float64) string {
	return fmt.Sprintf("%.2f%%", value)
}
