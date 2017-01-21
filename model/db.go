package model

import (
	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL Driver silly linter
	"gopkg.in/mgutz/dat.v1"
	"gopkg.in/mgutz/dat.v1/sqlx-runner"
	"log"
	"time"
)

var (
	db *runner.DB
)

// InitDB initializes the database pool and stores it in the model package scope.
func InitDB(ds string) {
	var err error
	conn, err := sql.Open("postgres", ds)

	if err != nil {
		log.Panic(err)
	}

	runner.MustPing(conn)

	// TODO: Have this be a parameter in the configuration
	conn.SetMaxIdleConns(4)
	conn.SetMaxOpenConns(16)

	dat.Strict = false

	// Log error if query takes over 25 milliseconds
	runner.LogQueriesThreshold = 25 * time.Millisecond

	db = runner.NewDB(conn, "postgres")
}
