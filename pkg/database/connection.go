package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	DB *sql.DB
}

func New(url string) (*DB, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	return &DB{DB: db}, db.Ping()
}

func (db *DB) Close() error {
	return db.DB.Close()
}