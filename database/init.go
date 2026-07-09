package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init(connectionString string) error {
	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	return nil
}
