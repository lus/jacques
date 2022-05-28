package discord

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

// Service represents the service handling the interaction with the Discord gateway
type Service struct {
	BotToken string

	session *discordgo.Session
}

// Start starts the Discord service
func (service *Service) Start() error {
	if service.session != nil {
		return errors.New("session already initialized")
	}

	session, err := discordgo.New("Bot " + service.BotToken)
	if err != nil {
		return err
	}
	if err := session.Open(); err != nil {
		return err
	}
	service.session = session
	return nil
}

// Stop stops the Discord service
func (service *Service) Stop() error {
	if err := service.session.Close(); err != nil {
		return err
	}
	service.session = nil
	return nil
}
