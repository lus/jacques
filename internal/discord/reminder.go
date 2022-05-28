package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zekrotja/ken"
)

type ReminderCommand struct{}

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
	return ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "soon:tm:",
		},
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
