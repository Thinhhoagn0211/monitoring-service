package main

import (
	"database/sql"
	"training/file-search/api"
	"training/file-search/util"

	db "training/file-search/db/sqlc"

	"github.com/rs/zerolog/log"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to db")
	}
	defer conn.Close()
	store := db.NewStore(conn)
	runGinServer(config, store)
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot start server")
	}
}
