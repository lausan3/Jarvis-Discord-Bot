package commands

import (
	"jarvis/utils/responses"

	"github.com/aws/aws-lambda-go/events"
	"github.com/bwmarrin/discordgo"
)

func EchoHandler(appToken string, interaction *discordgo.InteractionCreate) (events.APIGatewayV2HTTPResponse, error) {
	appCommand := interaction.ApplicationCommandData()
	var message = appCommand.GetOption("message").StringValue()

	if message == "" {
		message = "You didn't provide a message!"
	}

	response := discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	}

	return responses.NewSuccessfulBotResponse(appToken, map[string]any{
		"type": response.Type,
		"data": response.Data,
	})
}

var ApplicationCommandEcho discordgo.ApplicationCommand = discordgo.ApplicationCommand{
	Name:        "echo",
	Type:        discordgo.ChatApplicationCommand,
	Description: "Echoes your message.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        "message",
			Description: "The message to echo",
			Type:        3,
			Required:    true,
		},
	},
}
