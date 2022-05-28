package discord

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
	"github.com/lus/jacques/internal/reminder"
	"github.com/rs/zerolog/log"
	"github.com/zekrotja/ken"
	"strconv"
	"strings"
	"time"
)

type ReminderCommand struct {
	Repository reminder.Repository
}

var _ ken.SlashCommand = (*ReminderCommand)(nil)

func (cmd *ReminderCommand) Name() string {
	return "reminder"
}

func (cmd *ReminderCommand) Description() string {
	return "Create and manage things Jacques should remind you of"
}

func (cmd *ReminderCommand) Version() string {
	return "1.0.0"
}

func (cmd *ReminderCommand) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "list",
			Description: "List your pending reminders",
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "create",
			Description: "Create a new reminder",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "when",
					Description: "When the reminder should fire (e.g. '2h30m')",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "description",
					Description: "What the reminder is about",
					Required:    true,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "delete",
			Description: "Delete one or multiple of your reminders",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "id",
					Description: "The ID of the reminder to delete; '*' to delete all of your reminders",
					Required:    true,
				},
			},
		},
	}
}

func (cmd *ReminderCommand) Run(ctx *ken.Ctx) (err error) {
	return ctx.HandleSubCommands(
		ken.SubCommandHandler{
			Name: "list",
			Run:  cmd.listCommand,
		},
		ken.SubCommandHandler{
			Name: "create",
			Run:  cmd.createCommand,
		},
		ken.SubCommandHandler{
			Name: "delete",
			Run:  cmd.deleteCommand,
		},
	)
}

func (cmd *ReminderCommand) listCommand(ctx *ken.SubCommandCtx) error {
	ctx.SetEphemeral(true)
	return ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "soon:tm:",
		},
	})
}

func (cmd *ReminderCommand) createCommand(ctx *ken.SubCommandCtx) error {
	ctx.SetEphemeral(true)

	dur, err := time.ParseDuration(ctx.Options().GetByName("when").StringValue())
	if err != nil {
		return ctx.RespondError("The given duration is malformed.", "Invalid Duration")
	}
	description := strings.ReplaceAll(ctx.Options().GetByName("description").StringValue(), "`", "'")
	userID, _ := strconv.ParseInt(ctx.Event.Member.User.ID, 10, 64)
	channelID, _ := strconv.ParseInt(ctx.Event.ChannelID, 10, 64)

	rem, err := cmd.Repository.Create(context.Background(), &reminder.Create{
		UserID:      userID,
		ChannelID:   channelID,
		Description: description,
		Delta:       dur,
	})
	if err != nil {
		log.Error().Err(err).Msg("could not create reminder")
		return ctx.RespondError("An error occurred while creating the reminder.", "Error")
	}

	return ctx.RespondEmbed(&discordgo.MessageEmbed{
		Title:       "Success",
		Description: fmt.Sprintf("Your reminder `%s` will fire in %s.", rem.ID, durafmt.Parse(dur)),
		Color:       0x00ff00,
	})
}

func (cmd *ReminderCommand) deleteCommand(ctx *ken.SubCommandCtx) error {
	ctx.SetEphemeral(true)
	return ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "soon:tm:",
		},
	})
}

func (service *Service) fireReminder(rem *reminder.Reminder) {
	_, err := service.session.ChannelMessageSendComplex(strconv.FormatInt(rem.ChannelID, 10), &discordgo.MessageSend{
		Content: "<@" + strconv.FormatInt(rem.UserID, 10) + ">",
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Reminder",
				Description: "```" + rem.Description + "```",
				Color:       0xffff00,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("could not send reminder message")
	}
}
