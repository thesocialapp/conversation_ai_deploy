package main

import (
	"os"

	server "github.com/thesocialapp/conversation-ai/go/server"
	util "github.com/thesocialapp/conversation-ai/go/util"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	/// Find the env file in the parent dir and load it
	config, err := util.LoadConfig(".env")
	if err != nil {
		log.Fatal().Msgf("cannot load config %s", err.Error())
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msgf("Starting server at %v", config.HttpServerAddress)
	setupGinServer(config)
}

// Initializes Gin server and loads up config from .env
func setupGinServer(config util.Config) {
	server, err := server.NewServer(config)
	if err != nil {
		log.Error().Err(err).Msg("cannot create server")
	}
	err = server.StartServer()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start server")
	}
}
