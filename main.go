package main

import (
	api "conversational_ai/api/server"
	"conversational_ai/util"

	"github.com/rs/zerolog/log"
)

func main() {
	config, err := util.LoadConfig(".env")
	if err != nil {
		log.Fatal().Msgf("cannot load config %s", err.Error())
	}

	log.Info().Msgf("Starting server at %v", config.HttpServerAddress)
	setupGinServer(config)
}

// / Initializes the Gin server
func setupGinServer(config util.Config) {
	server, err := api.NewServer(config)
	if err != nil {
		log.Error().Err(err).Msg("cannot create server")
	}
	err = server.StartServer()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}
