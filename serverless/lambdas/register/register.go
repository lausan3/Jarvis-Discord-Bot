package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"jarvis/models"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

func handler(ctx context.Context) error {
	appID := os.Getenv("DISCORD_APP_ID")
	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	apiUrl := os.Getenv("DISCORD_API_URL")

	if appID == "" || botToken == "" || apiUrl == "" {
		return fmt.Errorf("missing required environment variables")
	}

	commands := []models.Command{
		{
			Name:        "ping",
			Type:        1,
			Description: "Replies with pong!",
		},
		{
			Name:        "echo",
			Type:        1,
			Description: "Echoes your message.",
			Options: []models.CommandOption{
				{
					Name:        "message",
					Description: "The message to echo",
					Type:        3,
					Required:    true,
				},
			},
		},
		{
			Name:        "help",
			Type:        1,
			Description: "Provides help information about the bot.",
		},
	}

	body, err := json.Marshal(commands)
	if err != nil {
		return fmt.Errorf("failed to serialize commands: %v", err)
	}

	url := fmt.Sprintf("%s/applications/%s/commands", apiUrl, appID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bot "+botToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord API call failed: %v\n%s", resp.Status, string(respBody))
	}

	fmt.Println("✅ Successfully registered commands with Discord!")

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %v", err)
	}

	stackName := "Register"
	cf := cloudformation.NewFromConfig(cfg)
	_, err = cf.DeleteStack(ctx, &cloudformation.DeleteStackInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete stack: %v", err)
	}

	fmt.Printf("🧨 Self-deletion of stack '%s' requested.\n", stackName)
	return nil
}

func main() {
	lambda.Start(handler)
}
