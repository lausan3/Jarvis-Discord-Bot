package main

import (
	"context"
	"encoding/json"
	chat_commands "jarvis/internal/commands/chat"
	message_commands "jarvis/internal/commands/message"
	"jarvis/utils/middleware"
	"jarvis/utils/responses"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Verify request signature and timestamp
	signature := request.Headers["x-signature-ed25519"]
	timestamp := request.Headers["x-signature-timestamp"]
	pubkey := os.Getenv("DISCORD_PUBLIC_KEY")
	token := os.Getenv("DISCORD_BOT_TOKEN")
	openaiKey := os.Getenv("OPENAI_API_KEY")
	body := request.Body

	logrus.Infof("Request: %v", request)

	if !middleware.ValidateDiscordSecurityHeaders(body, signature, timestamp, pubkey) {
		return responses.NewCustomErrorResponse(401, "Invalid request signature", nil)
	}

	// Parse interaction
	var interaction discordgo.InteractionCreate

	err := json.Unmarshal([]byte(body), &interaction)
	if err != nil {
		resp := map[string]any{
			"detail":       err.Error(),
			"bodyReceived": request.Body,
		}

		return responses.NewServerErrorGatewayResponse("Failed to parse interaction", resp)
	}

	// Respond to regular pings by Discord
	if interaction.Type == discordgo.InteractionPing {
		return responses.NewSuccessfulBotResponse(token, map[string]any{
			"type": discordgo.InteractionResponsePong,
		})
	}

	if interaction.Type != discordgo.InteractionApplicationCommand {
		return responses.NewClientErrorGatewayResponse("Unsupported interaction type", map[string]any{"bodyReceived": request.Body})
	}

	appCommand := interaction.ApplicationCommandData()

	commandType := appCommand.CommandType
	commandName := appCommand.Name

	// NOTE: Add new commands here
	switch commandType {
	case discordgo.ChatApplicationCommand:
		switch commandName {
		case "echo":
			return chat_commands.EchoCommandHandler(token, &interaction)
		}
	case discordgo.MessageApplicationCommand:
		switch commandName {
		case "summarize":
			messageId := appCommand.TargetID

			return message_commands.SummarizeMessageInteractionCommandHandler(token, openaiKey, messageId, &interaction)
		}
	}

	return responses.NewClientErrorGatewayResponse("Unknown command", map[string]any{"bodyReceived": request.Body})
}

func main() {
	lambda.Start(handler)
}
