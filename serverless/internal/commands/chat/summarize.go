package chat_commands

import (
	bot_utils "jarvis/utils/bot"
	bot_responses "jarvis/utils/bot/responses"
	"jarvis/utils/responses"

	"github.com/aws/aws-lambda-go/events"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// TODO: Find a way to make this more seamless than having to copy a message id.
var ChatApplicationCommandSummarize discordgo.ApplicationCommand = discordgo.ApplicationCommand{
	Name:        "summarize",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Summarizes recent chat messages.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "before_message",
			Description: "The message ID to summarize before",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	},
}

func SummarizeChatInteractionCommandHandler(c *discordgo.Session, botToken string, oaiKey string, interaction *discordgo.InteractionCreate) (events.APIGatewayV2HTTPResponse, error) {
	appCommand := interaction.ApplicationCommandData()
	var beforeMessageId = appCommand.GetOption("before_message").StringValue()

	if beforeMessageId == "" {
		logrus.Errorf("No message ID provided for summarize command")

		c.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: "You must provide a message ID to summarize before",
		})

		return responses.NewClientErrorGatewayResponse("You must provide a message ID to summarize before", nil)
	}

	member := interaction.Member

	if member == nil {
		logrus.Errorf("No member information in interaction. Did you use this handler in a user setting?")

		bot_responses.RespondUnexpectedError(c, interaction.Interaction)

		return responses.NewClientErrorGatewayResponse("This command can only be used in a server (guild) context", nil)
	}

	logrus.Infof("Received summarize command from user %s", member.User.Username)

	// Acknowledge the command to avoid the interaction timing out
	acknowledge := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Summarizing messages...",
		},
	}

	if err := c.InteractionRespond(interaction.Interaction, &acknowledge); err != nil {
		logrus.Errorf("Failed to send interaction response: %v", err)

		bot_responses.RespondUnexpectedError(c, interaction.Interaction)

		return responses.NewClientErrorGatewayResponse("Failed to send interaction response", map[string]any{"detail": err.Error()})
	}

	content, err := bot_utils.SummarizeBefore(botToken, oaiKey, beforeMessageId, interaction)
	if err != nil {
		logrus.Errorf("Error in summarize command: %v", err)

		c.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
			Content: "Failed to process summarize command. Please let Ant know. He may have run out of money.",
		})

		return responses.NewServerErrorGatewayResponse("Failed to process summarize command", map[string]any{"detail": err.Error()})
	}

	if _, err := c.FollowupMessageCreate(interaction.Interaction, true, &discordgo.WebhookParams{
		Content: content,
	}); err != nil {
		logrus.Errorf("Failed to send followup message: %v", err)

		bot_responses.RespondUnexpectedError(c, interaction.Interaction)

		return responses.NewServerErrorGatewayResponse("Failed to send interaction response", map[string]any{"detail": err.Error()})
	}

	return responses.NewSuccessGatewayResponse("Summarize command processed successfully", nil)
}
