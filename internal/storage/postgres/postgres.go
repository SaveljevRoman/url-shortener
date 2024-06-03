package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
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

	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS url (id SERIAL PRIMARY KEY, alias TEXT NOT NULL UNIQUE, url TEXT NOT NULL);`)
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return nil, fmt.Errorf("ошибка отката транзакции: %w", err)
		}
		return nil, fmt.Errorf("ошибка подготовки запроса создания таблицы url: %w", err)
	}

	_, err = tx.Exec(`CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);`)
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return nil, fmt.Errorf("ошибка отката транзакции: %w", err)
		}
		return nil, fmt.Errorf("ошибка создания индекса idx_alias: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return nil, fmt.Errorf("ошибка отката транзакции: %w", err)
		}
		return nil, fmt.Errorf("ошибка коммита транзакции: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	var id int64
	tx, err := s.db.Begin()
	if err != nil {
		return id, fmt.Errorf("ошибка старта транзакции сохранения url: %w", err)
	}

	query := `INSERT INTO url (alias, url) VALUES ($1, $2) RETURNING id`
	err = tx.QueryRow(query, alias, urlToSave).Scan(&id)
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return id, fmt.Errorf("ошибка отката транзакции: %v, ошибка вставки: %w", errRollback, err)
		}
		return id, fmt.Errorf("ошибка вставки и получения id: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return id, fmt.Errorf("ошибка коммита транзакции: %w", err)
	}

	return id, nil
}
