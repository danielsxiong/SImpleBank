package main

import (
	"danielsxiong/simplebank/api"
	db "danielsxiong/simplebank/db/sqlc"
	"danielsxiong/simplebank/util"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"

	"database/sql"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server")
	}
}
