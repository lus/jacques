package discord

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/lus/jacques/internal/reminder"
	"github.com/lus/jacques/internal/storage"
	"github.com/zekrotja/ken"
)

// Service represents the service handling the interaction with the Discord gateway
type Service struct {
	BotToken        string
	Storage         storage.Driver
	ReminderWatcher *reminder.Watcher

	session  *discordgo.Session
	commands *ken.Ken
}

// Start starts the Discord service
func (service *Service) Start() error {
	if service.session != nil {
		return errors.New("session already initialized")
	}

	// Create a new Discord session
	session, err := discordgo.New("Bot " + service.BotToken)
	if err != nil {
		return err
	}

	// Register all commands
	commands, err := ken.New(session)
	if err != nil {
		return err
	}
	err = commands.RegisterCommands(
		&ReminderCommand{
			Repository: &reminder.WatcherResettingRepository{
				Wrapping: service.Storage.Reminders(),
				Watcher:  service.ReminderWatcher,
			},
		},
	)
	if err != nil {
		return err
	}

	// Open the Discord session
	if err := session.Open(); err != nil {
		return err
	}

	service.session = session
	service.commands = commands
	service.ReminderWatcher.Subscribe(service.fireReminder)

	return nil
}

// Stop stops the Discord service
func (service *Service) Stop() error {
	if err := service.session.Close(); err != nil {
		return err
	}
	service.session = nil

	if err := service.commands.Unregister(); err != nil {
		return err
	}
	service.commands = nil

	return nil
}
