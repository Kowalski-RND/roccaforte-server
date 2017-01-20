package model

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL Driver silly linter
	"log"
)

var (
	db *sqlx.DB
)

// InitDB initializes the database pool and stores it in the model package scope.
func InitDB(ds string) {
	var err error
	db, err = sqlx.Connect("postgres", ds)

	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}
