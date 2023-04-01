package main

import (
	"danielsxiong/simplebank/api"
	db "danielsxiong/simplebank/db/sqlc"
	_ "github.com/lib/pq"

	"database/sql"
	"log"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:ab@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server")
	}
}
