package message_commands

import (
	bot_utils "jarvis/utils/bot"
	bot_responses "jarvis/utils/bot/responses"
	"jarvis/utils/responses"

	"github.com/aws/aws-lambda-go/events"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var MessageApplicationCommandSummarize = discordgo.ApplicationCommand{
	Name: "summarize",
	Type: discordgo.MessageApplicationCommand,
	IntegrationTypes: &[]discordgo.ApplicationIntegrationType{
		discordgo.ApplicationIntegrationUserInstall,
	},
}

func SummarizeMessageInteractionCommandHandler(c *discordgo.Session, botToken string, oaiKey string, beforeMessageId string, interaction *discordgo.InteractionCreate) (events.APIGatewayV2HTTPResponse, error) {
	member := interaction.Member

	if member == nil {
		logrus.Errorf("No member information in interaction. Did you use this handler in a user setting?")
		return responses.NewClientErrorGatewayResponse("This command can only be used in a server (guild) context", nil)
	}

	logrus.Infof("Received summarize command from user %s", member.User.Username)

	acknowledge := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
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

		bot_responses.RespondUnexpectedError(c, interaction.Interaction)

		return responses.NewServerErrorGatewayResponse("Failed to process summarize command", map[string]any{"detail": err.Error()})
	}

	c.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})

	return responses.NewSuccessGatewayResponse("Summarize command processed successfully", nil)
}
