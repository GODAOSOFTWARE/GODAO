// Структуры для работы с базой данных

package repository

import (
	"database/sql"
	"fmt"
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

// GetDB возвращает ссылку на базу данных
func GetDB() *sql.DB {
	return db
}

// Функция для получения названий всех таблиц в базе данных
func GetTableNames() ([]string, error) {
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tableNames = append(tableNames, name)
	}

	return tableNames, nil
}

// Функция для получения всех элементов из таблицы по её названию
func GetTableElements(tableName string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var elements []map[string]interface{}
	for rows.Next() {
		columnValues := make([]interface{}, len(columns))
		columnPointers := make([]interface{}, len(columns))
		for i := range columnValues {
			columnPointers[i] = &columnValues[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		element := make(map[string]interface{})
		for i, colName := range columns {
			val := columnPointers[i].(*interface{})
			element[colName] = *val
		}

		elements = append(elements, element)
	}

	return elements, nil
}
