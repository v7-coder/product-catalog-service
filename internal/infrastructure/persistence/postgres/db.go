package postgres

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func NewDB(dsn string) (*DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(50) // Максимум открытых соединений
	db.SetMaxIdleConns(25) // Простаивающие соединения
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
