package services

import (
	"dao_vote/internal/models"
	"dao_vote/internal/repository"
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
