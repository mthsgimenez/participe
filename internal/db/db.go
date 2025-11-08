package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectToDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("connect_to_db: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect_to_db: %w", err)
	}

	return db, nil
}

