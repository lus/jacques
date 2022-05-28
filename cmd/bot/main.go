package main

import (
	"github.com/lus/jacques/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

func main() {
	// Set up zerolog to use pretty printing
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out: os.Stderr,
	})
	log.Info().Msg("starting up...")

	// Load the configuration
	log.Info().Msg("loading configuration...")
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("could not load configuration")
	}
	_ = cfg

	// Wait for the application to be terminated
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
}
