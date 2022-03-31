package service

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	db "github.com/nemo984/money-app-api/db/sqlc"
)

var (
	testService Service
	dbDriver    = "postgres"
	dbSource    = "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
)

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	queries := db.New(conn)
	testService = New(queries)

	os.Exit(m.Run())
}
