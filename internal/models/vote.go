package models

// Vote представляет структуру для пользовательского голосования.
type Vote struct {
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

// VoteWithoutID представляет структуру для пользовательского голосования без VoterID.
type VoteWithoutID struct {
	Title       string `json:"title" validate:"required"`       // Заголовок голосования
	Subtitle    string `json:"subtitle" validate:"required"`    // Подзаголовок голосования
	Description string `json:"description" validate:"required"` // Описание предложения
	Voter       string `json:"voter" validate:"required"`       // Адрес кошелька, с которого было отправлено голосование
	Choice      string `json:"choice" validate:"required"`      // Выбранный вариант голосования ("За" или "Против")
}

// DAOTeamApiResponse представляет ответ от API результатов голосования команды DAO.
type DAOTeamApiResponse struct {
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

// DAOTeamVoteResultsResponse представляет обработанные результаты голосования команды DAO.
type DAOTeamVoteResultsResponse struct {
	DAOMembers        int           `json:"dao_members"`
	TotalTransactions int           `json:"total_transactions"`
	VotedMembers      int           `json:"voted_members"`
	Turnout           string        `json:"turnout"`
	VotesFor          string        `json:"votes_for"`
	VotesAgainst      string        `json:"votes_against"`
	VotingStatus      string        `json:"voting_status"`
	Resolution        string        `json:"resolution"`
	ValidTransactions []Transaction `json:"transactions"`
	RejectedTxs       []Transaction `json:"rejected_transactions"`
	NullVotePowerTxs  []Transaction `json:"null_vote_power_transactions"`
	InvalidMessageTxs []Transaction `json:"invalid_message_transactions"`
}

// UserVote представляет структуру для голосов пользователей.
type UserVote struct {
	VoterID   int    `json:"id"`         // Уникальный идентификатор голоса
	VoteID    int    `json:"vote_id"`    // VoterID голосования
	Voter     string `json:"voter"`      // Адрес кошелька голосующего
	Choice    string `json:"choice"`     // Выбранный вариант ("За" или "Против")
	VotePower int    `json:"vote_power"` // Сила голоса
}
