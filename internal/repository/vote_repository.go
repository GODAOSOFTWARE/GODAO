package repository

import (
	"dao_vote/internal/models"
	"errors"
	"sync"
)

// Карта силы голосов для различных кошельков
var voteMap = map[string]int{
	"d016nnqrut83vd0p4afp6546rma6g5d8aqy6t7cfp": 1287398,
	"d012z4fvdrdlwp54lhqgay2j37hyk93z8ce44w5sx": 432000,
	// (остальные элементы карты)
	// ...
	"d010uuxz0wnc9x9zhd80x9kp9urghke9n4m6espn6": 32589,
}

// totalVoices представляет общее количество голосов
var totalVoices = 12876983

// GetVoteStrength возвращает силу голоса для указанного кошелька
func GetVoteStrength(from string) (int, error) {
	strength, exists := voteMap[from]
	if !exists {
		return 0, errors.New("сторонний голос")
	}
	return strength, nil
}

// GetTotalVoices возвращает общее количество голосов
func GetTotalVoices() int {
	return totalVoices
}

// GetVoteMap возвращает карту всех голосов
func GetVoteMap() map[string]int {
	return voteMap
}

// Временное хранилище для пользовательских голосований
var (
	votes  = make(map[int]models.Vote) // Словарь для хранения голосований
	mu     sync.Mutex                  // Мьютекс для управления доступом к хранилищу
	nextID = 1                         // ID для следующего голосования
)

// SaveVote сохраняет новое пользовательское голосование
func SaveVote(vote models.Vote) (int, error) {
	mu.Lock()         // Захватываем мьютекс
	defer mu.Unlock() // Освобождаем мьютекс после выполнения функции
	vote.ID = nextID  // Присваиваем новое ID голосованию
	votes[nextID] = vote
	nextID++ // Увеличиваем ID для следующего голосования
	return vote.ID, nil
}

// GetVoteByID возвращает пользовательское голосование по его ID
func GetVoteByID(id int) (models.Vote, error) {
	mu.Lock()
	defer mu.Unlock()
	vote, exists := votes[id]
	if !exists {
		return models.Vote{}, errors.New("голосование не найдено")
	}
	return vote, nil
}

// DeleteVote удаляет пользовательское голосование по его ID
func DeleteVote(id int) error {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := votes[id]; !exists {
		return errors.New("голосование не найдено")
	}
	delete(votes, id)
	return nil
}

// Временное хранилище для голосов пользователей
var (
	userVotes  = make(map[int][]models.UserVote)
	voteMu     sync.Mutex
	userVoteID = 1
)

// AddUserVote сохраняет новый голос пользователя
func AddUserVote(vote models.UserVote) (int, error) {
	voteMu.Lock()
	defer voteMu.Unlock()

	vote.VoterID = userVoteID
	userVotes[vote.VoteID] = append(userVotes[vote.VoteID], vote)
	userVoteID++
	return vote.VoterID, nil
}

// GetUserVotes возвращает все голоса пользователей для указанного голосования
func GetUserVotes(voteID int) ([]models.UserVote, error) {
	voteMu.Lock()
	defer voteMu.Unlock()

	votes, exists := userVotes[voteID]
	if !exists {
		return nil, errors.New("No user votes found for this vote")
	}
	return votes, nil
}
