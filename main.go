package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/phillipngn/bankie/api"
	db "github.com/phillipngn/bankie/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/bankie?sslmode=disable"
	address  = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect database")
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(address)

	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
