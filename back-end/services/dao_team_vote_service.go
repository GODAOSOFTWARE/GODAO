// Package services Логика для корпоративного голосования
package services

import (
	"dao_vote/back-end/models"
	"dao_vote/back-end/repository"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	requiredMajority = 51  // Требуемое большинство для принятия решения
	percentFactor    = 100 // Фактор для расчета процентов
)

// / Функция для получения количества записей в таблице vote_strength
func getVoteStrengthCount() (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM vote_strength"
	db := repository.GetDB()
	log.Println("Executing query to get count of vote_strength")
	if err := db.QueryRow(query).Scan(&count); err != nil {
		return 0, fmt.Errorf("error querying vote_strength count: %v", err)
	}
	log.Printf("Count of vote_strength: %d\n", count)
	return count, nil
}

// FetchDAOTeamVoteResults - функция для получения результатов голосования по адресу кошелька DAO
func FetchDAOTeamVoteResults(walletAddress string, offset int) (models.DAOTeamApiResponse, error) {
	log.Printf("Fetching DAO Team Vote Results for wallet: %s with offset: %d\n", walletAddress, offset)

	// Получаем количество членов ДАО
	limit, err := getVoteStrengthCount()
	if err != nil {
		log.Printf("Error getting DAO members count: %v\n", err)
		return models.DAOTeamApiResponse{}, err
	}
	log.Printf("Limit set to: %d\n", limit)

	// Формируем URL для запроса к API
	apiURL := fmt.Sprintf("https://mainnet-explorer-api.decimalchain.com/api/address/%s/txs?limit=%d&offset=%d", walletAddress, limit, offset)
	log.Printf("API URL: %s\n", apiURL)

	// Выполняем GET-запрос к API
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("Error making request: %v\n", err)
		return models.DAOTeamApiResponse{}, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем, что ответ от сервера успешен
	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-200 response code: %d\n", resp.StatusCode)
		return models.DAOTeamApiResponse{}, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	// Читаем тело ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return models.DAOTeamApiResponse{}, fmt.Errorf("error reading response body: %v", err)
	}
	log.Printf("Response body: %s\n", string(body))

	// Парсим JSON-ответ в структуру DAOTeamApiResponse
	var apiResponse models.DAOTeamApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Error unmarshalling response body: %v\n", err)
		return models.DAOTeamApiResponse{}, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	// Обновляем силу голосов для каждой транзакции в ответе
	for i, result := range apiResponse.Result.Txs {
		log.Printf("Processing transaction from: %s\n", result.From)
		votePower, err := repository.GetVoteStrength(result.From)
		if err != nil {
			log.Printf("Error getting vote strength for %s: %v\n", result.From, err)
		}
		apiResponse.Result.Txs[i].VotePower = votePower
		apiResponse.Result.Txs[i].Hash = result.Hash
		log.Printf("Updated transaction: %+v\n", apiResponse.Result.Txs[i])
	}

	log.Printf("Final API response: %+v\n", apiResponse)
	return apiResponse, nil
}

// PrepareDAOTeamVoteResults - функция для подготовки результатов голосования команды DAO
func PrepareDAOTeamVoteResults(apiResponse models.DAOTeamApiResponse) models.DAOTeamVoteResultsResponse {
	// Инициализируем списки для различных категорий транзакций
	validTxs := []models.Transaction{}
	duplicateTxs := []models.Transaction{}
	nullVotePowerTxs := []models.Transaction{}
	invalidTxs := []models.Transaction{}
	votesFor := []models.Transaction{}
	votesAgainst := []models.Transaction{}
	uniqueVoters := make(map[string]bool) // Карта уникальных голосующих
	totalTransactions := len(apiResponse.Result.Txs)

	// Обрабатываем каждую транзакцию
	for _, result := range apiResponse.Result.Txs {
		// Приводим сообщение к нижнему регистру и удаляем лишние пробелы и кавычки
		message := strings.TrimSpace(strings.ToLower(result.Message))
		message = strings.Trim(message, `\"`)

		// Логируем детали транзакции
		log.Printf("Processing transaction from: %s, message: %s, vote power: %d, hash: %s", result.From, message, result.VotePower, result.Hash)

		// Проверка на нулевую силу голоса
		if result.VotePower == 0 {
			nullVotePowerTxs = append(nullVotePowerTxs, result)
			log.Printf("Transaction with zero vote power: %s", result.Hash)
			continue
		}

		// Проверка на пустое сообщение
		if message == "" {
			invalidTxs = append(invalidTxs, result)
			log.Printf("Invalid transaction with empty message: %s", result.Hash)
			continue
		}

		// Проверка на дублирующие транзакции
		if uniqueVoters[result.From] {
			duplicateTxs = append(duplicateTxs, result)
			log.Printf("Duplicate transaction: %s", result.Hash)
			continue
		}

		// Классификация транзакции на "За" или "Против"
		switch message {
		case "да", "дa", "д", "за", "зa", "z":
			votesFor = append(votesFor, result)
			validTxs = append(validTxs, result)
			log.Printf("Vote for transaction: %s", result.Hash)
		case "нет", "н", "против":
			votesAgainst = append(votesAgainst, result)
			validTxs = append(validTxs, result)
			log.Printf("Vote against transaction: %s", result.Hash)
		default:
			invalidTxs = append(invalidTxs, result)
			log.Printf("Invalid transaction with unknown message: %s", result.Hash)
			continue
		}

		// Помечаем голосующего как уникального
		uniqueVoters[result.From] = true
	}

	// Вычисляем силу голосов за и против
	strengthFor := calculateStrength(votesFor)
	strengthAgainst := calculateStrength(votesAgainst)
	totalVoices := repository.GetTotalVoices()
	percentFor := calculatePercentage(strengthFor, totalVoices)
	percentAgainst := calculatePercentage(strengthAgainst, totalVoices)

	// Определяем статус голосования и резолюцию
	status := "Активно"
	resolution := "Решение не принято"
	if percentFor >= requiredMajority {
		status = "Завершено"
		resolution = "Принять изменения"
	} else if percentAgainst >= requiredMajority {
		status = "Завершено"
		resolution = "Отклонить изменения"
	}

	// Получаем количество членов ДАО из базы данных
	daoMembers, err := getVoteStrengthCount()
	if err != nil {
		log.Printf("Error getting DAO members count: %v\n", err)
		daoMembers = 0 // Устанавливаем 0, если возникла ошибка
	}

	// Логируем итоговые результаты
	log.Printf("Voting completed with %d DAO members, %d voted members", daoMembers, len(uniqueVoters))
	log.Printf("Votes for: %d/%d (%.2f%%), Votes against: %d (%.2f%%)", strengthFor, totalVoices, percentFor, strengthAgainst, percentAgainst)
	log.Printf("Voting status: %s, Resolution: %s", status, resolution)

	// Возвращаем результаты голосования
	return models.DAOTeamVoteResultsResponse{
		DAOMembers:        daoMembers,
		VotedMembers:      len(uniqueVoters),
		Turnout:           formatPercentage(float64(len(uniqueVoters)) / float64(daoMembers) * percentFactor),
		VotesFor:          fmt.Sprintf("%d/%d (%.2f%%)", strengthFor, totalVoices, percentFor),
		VotesAgainst:      fmt.Sprintf("%d (%.2f%%)", strengthAgainst, percentAgainst),
		VotingStatus:      status,
		Resolution:        resolution,
		ValidTransactions: validTxs,          // Валидные транзакции
		TotalTransactions: totalTransactions, // Общее количество транзакций
		RejectedTxs:       duplicateTxs,      // Задвоенные транзакции
		NullVotePowerTxs:  nullVotePowerTxs,  // Транзакции с нулевой силой голоса
		InvalidMessageTxs: invalidTxs,        // Транзакции с некорректным сообщением
	}
}

// calculateStrength - функция для вычисления общей силы голосов из списка транзакций
func calculateStrength(votes []models.Transaction) int {
	strength := 0
	// Суммируем силу голосов для каждой транзакции
	for _, vote := range votes {
		if voteStrength, err := repository.GetVoteStrength(vote.From); err == nil {
			strength += voteStrength
		} else {
			fmt.Printf("Error getting vote strength for %s: %v\n", vote.From, err)
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
