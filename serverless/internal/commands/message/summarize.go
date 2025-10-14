package message_commands

import (
	bot_utils "jarvis/utils/bot"
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

func SummarizeMessageInteractionCommandHandler(botToken string, oaiKey string, beforeMessageId string, interaction *discordgo.InteractionCreate) (events.APIGatewayV2HTTPResponse, error) {
	member := interaction.Member

	if member == nil {
		logrus.Errorf("No member information in interaction. Did you use this handler in a user setting?")
		return responses.NewClientErrorGatewayResponse("This command can only be used in a server (guild) context", nil)
	}

	logrus.Infof("Received summarize command from user %s", member.User.Username)

	content, err := bot_utils.SummarizeBefore(botToken, oaiKey, beforeMessageId, interaction)
	if err != nil {
		logrus.Errorf("Error in summarize command: %v", err)
		return responses.NewServerErrorGatewayResponse("Failed to process summarize command", map[string]any{"detail": err.Error()})
	}

	return responses.NewSuccessfulBotResponse(botToken, map[string]any{
		"type": discordgo.InteractionResponseChannelMessageWithSource,
		"data": map[string]any{
			"content": content,
		},
	})
}
