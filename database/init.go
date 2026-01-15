package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("postgres", os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatal("database connection failed", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("database unreachable: ", err)
	}
}
