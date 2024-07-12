// Логика для голосования пользователей и обработки результатов голосования

package services

import (
	"dao_vote/back-end/models"
	"dao_vote/back-end/repository"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// CreateVote создает новое пользовательское голосование и возвращает его ID.
func CreateVote(vote models.Vote) (int, error) {
	return repository.SaveVote(vote)
}

// GetVote получает пользовательское голосование по ID.
func GetVote(id int) (models.Vote, error) {
	return repository.GetVoteByID(id)
}

// DeleteVote удаляет пользовательское голосование по ID.
func DeleteVote(id int) error {
	return repository.DeleteVote(id)
}

// GetVoteStrength возвращает силу голоса для указанного кошелька.
func GetVoteStrength(from string) (int, error) {
	return repository.GetVoteStrength(from)
}

// AddUserVote добавляет новый голос пользователя и возвращает его ID.
func AddUserVote(vote models.UserVote) (int, error) {
	return repository.AddUserVote(vote)
}

// GetUserVotes получает все голоса пользователей для указанного голосования.
func GetUserVotes(voteID int) ([]models.UserVote, error) {
	return repository.GetUserVotes(voteID)
}

// FetchUserVotes получает результаты голосования по ID голосования
func FetchUserVotes(voteID int) (models.DAOTeamVoteResultsResponse, error) {
	// Получаем голосование по ID
	vote, err := repository.GetVoteByID(voteID)
	if err != nil {
		return models.DAOTeamVoteResultsResponse{}, fmt.Errorf("failed to get vote by ID: %v", err)
	}

	// Логируем адрес кошелька для голосования
	logrus.Infof("Parsing wallet address for vote ID %d: %s", voteID, vote.WalletAddress)

	// Формируем URL для запроса транзакций кошелька
	apiURL := fmt.Sprintf("https://mainnet-explorer-api.decimalchain.com/api/address/%s/txs", vote.WalletAddress)
	resp, err := http.Get(apiURL)
	if err != nil {
		return models.DAOTeamVoteResultsResponse{}, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.DAOTeamVoteResultsResponse{}, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.DAOTeamVoteResultsResponse{}, fmt.Errorf("error reading response body: %v", err)
	}

	// Парсим ответ API
	var apiResponse models.DAOTeamApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return models.DAOTeamVoteResultsResponse{}, fmt.Errorf("error unmarshalling response body: %v", err)
	}

	// Обновляем силу голосов и хэши для каждой транзакции
	for i, result := range apiResponse.Result.Txs {
		votePower, err := repository.GetVoteStrength(result.From)
		if err != nil {
			fmt.Printf("Error getting vote strength for %s: %v\n", result.From, err)
		}
		apiResponse.Result.Txs[i].VotePower = votePower
		apiResponse.Result.Txs[i].Hash = result.Hash
	}

	// Возвращаем обработанные результаты голосования
	return PrepareDAOTeamVoteResults(apiResponse), nil
}
