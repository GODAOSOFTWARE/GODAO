// repository/db.go

package repository

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// InitDB инициализирует базу данных
func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}

	// Создаем таблицу для голосований, если она не существует
	createVotesTable := `
    CREATE TABLE IF NOT EXISTS votes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT,
        subtitle TEXT,
        description TEXT,
        voter TEXT,
        choice TEXT,
        vote_power INTEGER,
        wallet_address TEXT
    );`
	if _, err := db.Exec(createVotesTable); err != nil {
		return err
	}

	// Создаем таблицу для голосов пользователей, если она не существует
	createUserVotesTable := `
    CREATE TABLE IF NOT EXISTS user_votes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        vote_id INTEGER,
        voter TEXT,
        choice TEXT,
        vote_power INTEGER
    );`
	if _, err := db.Exec(createUserVotesTable); err != nil {
		return err
	}

	// Создаем таблицу для силы голосов кошельков, если она не существует
	createVoteStrengthTable := `
    CREATE TABLE IF NOT EXISTS vote_strength (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        wallet_address TEXT UNIQUE,
        vote_power INTEGER
    );`
	if _, err := db.Exec(createVoteStrengthTable); err != nil {
		return err
	}

	return nil
}
