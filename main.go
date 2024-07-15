package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/chowchow-dev/bankie/api"
	db "github.com/chowchow-dev/bankie/db/sqlc"
	"github.com/chowchow-dev/bankie/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect database")
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
