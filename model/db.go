package model

import (
	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL Driver silly linter
	"log"
)

var db *sql.DB

// InitDB initializes the database pool and stores it in the model package scope.
func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)

	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}
