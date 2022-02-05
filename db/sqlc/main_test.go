package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries

const dbSource = "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"

func TestMain(m *testing.M) {
	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
