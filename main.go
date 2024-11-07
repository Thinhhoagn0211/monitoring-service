package main

import (
	"database/sql"
	"filesearch/api"
	"filesearch/util"

	db "filesearch/db"

	"github.com/rs/zerolog/log"
)

const ()

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}
	store := db.NewStore(conn)
	runGinServer(config, store)
}

func runGinServer(config util.Config) {
	server, err := api.NewServer(config)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot start server")
	}
}
