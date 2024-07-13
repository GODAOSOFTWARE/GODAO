// Package models Структуры для голосования пользователей
package models

// VoteInfo представляет структуру для хранения пользовательского голосования.
type VoteInfo struct {
	ID             int    `json:"id"`                              // Уникальный идентификатор голосования
	Title          string `json:"title" validate:"required"`       // Заголовок голосования
	Subtitle       string `json:"subtitle" validate:"required"`    // Подзаголовок голосования
	Description    string `json:"description" validate:"required"` // Описание предложения
	Voter          string `json:"voter" validate:"required"`       // Адрес кошелька, с которого было отправлено голосование
	Choice         string `json:"choice" validate:"required"`      // Выбранный вариант голосования ("За" или "Против")
	VotePower      int    `json:"vote_power"`                      // Сила голоса
	WalletAddress  string `json:"wallet_address"`                  // Адрес кошелька
	MnemonicPhrase string `json:"-"`                               // Мнемоническая фраза, скрыта в JSON-ответах
}

// NewVote представляет структуру для пользовательского голосования без VoterID.
type NewVote struct {
	Title       string `json:"title" validate:"required"`       // Заголовок голосования
	Subtitle    string `json:"subtitle" validate:"required"`    // Подзаголовок голосования
	Description string `json:"description" validate:"required"` // Описание предложения
	Voter       string `json:"voter" validate:"required"`       // Адрес кошелька, с которого было отправлено голосование
	Choice      string `json:"choice" validate:"required"`      // Выбранный вариант голосования ("За" или "Против")
}

// WithdrawOrderResponse представляет ответ от API результатов голосования команды DAO.
type WithdrawOrderResponse struct {
	Result struct {
		Txs []Transaction `json:"txs"`
	} `json:"result"`
}

// Transaction представляет одну транзакцию в результатах голосования команды DAO.
type Transaction struct {
	From      string `json:"from"`
	Message   string `json:"message"`
	VotePower int    `json:"vote_power"`
	Hash      string `json:"hash"` // Добавлено поле для хэша транзакции
}

// VoteResults представляет обработанные результаты голосования команды DAO.
type VoteResults struct {
	DAOMembers        int           `json:"dao_members"`
	TotalTransactions int           `json:"total_transactions"`
	VotedMembers      int           `json:"voted_members"`
	Turnout           string        `json:"turnout"`
	VotesFor          string        `json:"votes_for"`
	VotesAgainst      string        `json:"votes_against"`
	VotingStatus      string        `json:"voting_status"`
	Resolution        string        `json:"resolution"`
	ValidTransactions []Transaction `json:"valid_transactions"`
	RejectedTxs       []Transaction `json:"rejected_transactions"`
	NullVotePowerTxs  []Transaction `json:"null_vote_power_transactions"`
	InvalidMessageTxs []Transaction `json:"invalid_message_transactions"`
}

// UserVote представляет структуру для голосjdfybz пользователей.
type UserVote struct {
	VoterID   int    `json:"id"`         // Уникальный идентификатор голоса
	VoteID    int    `json:"vote_id"`    // VoterID голосования
	Voter     string `json:"voter"`      // Адрес кошелька голосующего
	Choice    string `json:"choice"`     // Выбранный вариант ("За" или "Против")
	VotePower int    `json:"vote_power"` // Сила голоса
}
