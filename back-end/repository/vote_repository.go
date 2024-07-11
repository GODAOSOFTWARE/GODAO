package repository

import (
	"dao_vote/back-end/models"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

// Глобальная переменная для базы данных
var db *sql.DB

// Карта силы голосов для различных кошельков
var voteMap = map[string]int{
	// (данные карты)
	"d016nnqrut83vd0p4afp6546rma6g5d8aqy6t7cfp": 1287398,
	// (остальные данные)
	"d01p55v08ld8yc0my72ccpsztv7auyxn2tden6yvw": 1000000,
}

// GetVoteStrength возвращает силу голоса для указанного кошелька
func GetVoteStrength(from string) (int, error) {
	strength, exists := voteMap[from]
	if !exists {
		return 0, errors.New("сторонний голос")
	}
	return strength, nil
}

// GetVoteMap возвращает карту всех голосов
func GetVoteMap() map[string]int {
	return voteMap
}

// GetTotalVoices возвращает общее количество голосов
func GetTotalVoices() int {
	totalVoices := 0
	for _, v := range voteMap {
		totalVoices += v
	}
	return totalVoices
}

// SaveVote сохраняет новое пользовательское голосование
func SaveVote(vote models.Vote) (int, error) {
	result, err := db.Exec("INSERT INTO votes (title, subtitle, description, voter, choice, vote_power, wallet_address) VALUES (?, ?, ?, ?, ?, ?, ?)",
		vote.Title, vote.Subtitle, vote.Description, vote.Voter, vote.Choice, vote.VotePower, vote.WalletAddress)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetVoteByID возвращает пользовательское голосование по его ID
func GetVoteByID(id int) (models.Vote, error) {
	var vote models.Vote
	err := db.QueryRow("SELECT id, title, subtitle, description, voter, choice, vote_power, wallet_address FROM votes WHERE id = ?", id).Scan(
		&vote.ID, &vote.Title, &vote.Subtitle, &vote.Description, &vote.Voter, &vote.Choice, &vote.VotePower, &vote.WalletAddress)
	if err != nil {
		if err == sql.ErrNoRows {
			return vote, errors.New("голосование не найдено")
		}
		return vote, err
	}
	return vote, nil
}

// DeleteVote удаляет пользовательское голосование по его ID
func DeleteVote(id int) error {
	_, err := db.Exec("DELETE FROM votes WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

// AddUserVote сохраняет новый голос пользователя
func AddUserVote(vote models.UserVote) (int, error) {
	result, err := db.Exec("INSERT INTO user_votes (vote_id, voter, choice, vote_power) VALUES (?, ?, ?, ?)",
		vote.VoteID, vote.Voter, vote.Choice, vote.VotePower)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetUserVotes возвращает все голоса пользователей для указанного голосования
func GetUserVotes(voteID int) ([]models.UserVote, error) {
	rows, err := db.Query("SELECT id, vote_id, voter, choice, vote_power FROM user_votes WHERE vote_id = ?", voteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var votes []models.UserVote
	for rows.Next() {
		var vote models.UserVote
		if err := rows.Scan(&vote.VoterID, &vote.VoteID, &vote.Voter, &vote.Choice, &vote.VotePower); err != nil {
			return nil, err
		}
		votes = append(votes, vote)
	}
	return votes, nil
}
