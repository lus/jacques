package main

import (
	"context"
	"github.com/lus/jacques/internal/config"
	"github.com/lus/jacques/internal/discord"
	"github.com/lus/jacques/internal/reminder"
	"github.com/lus/jacques/internal/storage/postgres"
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

	// Establish the PostgreSQL connection
	log.Info().Msg("connecting to database...")
	driver := postgres.New(cfg.PostgresDSN)
	if err := driver.Initialize(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("could not connect to database")
	}
	defer func() {
		log.Info().Msg("disconnecting from database...")
		if err := driver.Close(); err != nil {
			log.Warn().Err(err).Msg("could not disconnect from database")
		}
	}()

	// Create and start the reminder watcher
	watcher := &reminder.Watcher{
		Repo: driver.Reminders(),
	}
	watcher.Start()
	defer watcher.Stop()

	// Start the Discord service
	log.Info().Msg("connecting to gateway...")
	dcService := &discord.Service{
		BotToken:        cfg.BotToken,
		Storage:         driver,
		ReminderWatcher: watcher,
	}
	if err := dcService.Start(); err != nil {
		log.Fatal().Err(err).Msg("could not connect to gateway")
	}
	defer func() {
		log.Info().Msg("disconnecting from gateway...")
		if err := dcService.Stop(); err != nil {
			log.Warn().Err(err).Msg("could not disconnect from gateway")
		}
	}()

	// Done!
	log.Info().Msg("done!")

	// Wait for the application to be terminated
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
}
