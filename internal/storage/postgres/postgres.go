package postgres

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(host, port, dbName, user, password string) (*Storage, error) {
	pgInfo := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, dbName, host, port,
	)
	db, err := sql.Open("postgres", pgInfo)
	if err != nil {
		return nil, fmt.Errorf("нет конекта с postgres: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("ошибка старта транзакции: %w", err)
	}

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS url (id INTEGER PRIMARY KEY, alias TEXT NOT NULL UNIQUE, url TEXT NOT NULL);`)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, fmt.Errorf("ошибка отката транзакции: %w", err)
		}
		return nil, fmt.Errorf("ошибка подготовки запроса создания таблицы url: %w", err)
	}

	_, err = tx.Exec(`CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);`)
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, fmt.Errorf("ошибка отката транзакции: %w", err)
		}
		return nil, fmt.Errorf("ошибка создания индекса idx_alias: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		if err = tx.Rollback(); err != nil {
			return nil, fmt.Errorf("ошибка отката транзакции: %w", err)
		}
		return nil, fmt.Errorf("ошибка коммита транзакции: %w", err)
	}

	return &Storage{db: db}, nil
}
