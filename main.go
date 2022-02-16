package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/handler"
	"github.com/nemo984/money-app-api/service"
)

var (
	dbDriver = os.Getenv("DB_DRIVER")
	dbSource = os.Getenv("DB_SOURCE")
	port = os.Getenv("PORT")
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	queries := db.New(conn)
	service := service.NewService(queries)
	s := handler.NewServer(service)

	log.Fatal(s.Start(fmt.Sprintf(":%s", port)))
}
