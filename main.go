package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	db "github.com/nemo984/money-app-api/db/sqlc"
	"github.com/nemo984/money-app-api/handler"
	"github.com/nemo984/money-app-api/notification"
	"github.com/nemo984/money-app-api/service"
)

var (
	dbDriver = os.Getenv("DB_DRIVER")
	dbSource = os.Getenv("DB_SOURCE")
	port     = os.Getenv("PORT")
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	queries := db.New(conn)
	service := service.New(queries)
	hub := notification.New()
	s := handler.New(service, hub)

	log.Fatal(s.Start(":" + port))
}
