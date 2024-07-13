// Package repository Хранилище для действий связанных с голосованиями
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
var voteMap = map[string]int{}

// GetVoteStrength возвращает силу голоса для указанного кошелька из базы данных
func GetVoteStrength(walletAddress string) (int, error) {
	var votePower int
	err := db.QueryRow("SELECT vote_power FROM vote_strength WHERE wallet_address = ?", walletAddress).Scan(&votePower)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("сторонний голос")
		}
		return 0, err
	}
	return votePower, nil
}

// AddWalletStrength добавляет новый адрес кошелька и силу голоса в базу данных
func AddWalletStrength(walletAddress string, votePower int) error {
	_, err := db.Exec("INSERT INTO vote_strength (wallet_address, vote_power) VALUES (?, ?)", walletAddress, votePower)
	return err
}

// DeleteWalletStrength удаляет адрес кошелька и силу голоса из базы данных
func DeleteWalletStrength(walletAddress string) error {
	_, err := db.Exec("DELETE FROM vote_strength WHERE wallet_address = ?", walletAddress)
	return err
}

// GetVoteMap возвращает карту всех голосов
func GetVoteMap() map[string]int {
	return voteMap
}

// GetTotalVoices возвращает общее количество голосов
func GetTotalVoices() int {
	totalVoices := 12000000
	for _, v := range voteMap {
		totalVoices += v
	}
	return totalVoices
}

// SaveVote сохраняет новое пользовательское голосование
func SaveVote(vote models.VoteInfo) (int, error) {
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
func GetVoteByID(id int) (models.VoteInfo, error) {
	var vote models.VoteInfo
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
