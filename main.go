package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/handler"
	"github.com/nemo984/money-app-api/service"
)

const dbSource = "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"

func main() {
	conn, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	queries := db.New(conn)
	service := service.NewService(queries)
	s := handler.NewServer(service)

	log.Fatal(s.Start(":8080"))
}
