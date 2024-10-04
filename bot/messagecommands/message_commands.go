package messagecommands

import (
	"context"
	"fmt"
	"main/infra/logger"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
)

// Send a test message as a response to "jarvis test"
func Test(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Hello!")
	if err != nil {
		logger.Errorf("Could not send test message response!")
		return
	}

	logger.Infof("Successfully sent message at %s", m.Timestamp)
}

func Summarize(s *discordgo.Session, m *discordgo.MessageCreate) {
	openAIToken := viper.GetString("OPENAI_TOKEN")
	if openAIToken == "" {
		logger.Errorf("Error during Jarvis Summarize: Couldn't find our OpenAI API token, perhaps it is not set?")
	}

	client := openai.NewClient(openAIToken)
	systemPrompt := `
		You will be given context in the form of a string representing messages in the format <USER NAME> said: <MESSAGE CONTENT>\n.
		Ignore any messages sent by YOU, Jarvis.

		Your task is to take this context and summarize what the conversation was about.
	`

	messagesArr, err := s.ChannelMessages(m.ChannelID, 5, m.ID, "", "")
	if err != nil {
		logger.Errorf("Error getting messages before message id %s: %s", m.ID, err.Error())
		return
	}

	messages := []string{}

	for _, message := range messagesArr {
		messages = append(messages, fmt.Sprintf("%s said: %s", message.Author.GlobalName, message.Content))
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: strings.Join(messages, "\n"),
				},
			},
		},
	)

	if err != nil {
		logger.Errorf("Chat Completion error: %v", err)
		return
	}

	if _, err = s.ChannelMessageSendReply(m.ChannelID, resp.Choices[0].Message.Content, m.Reference()); err != nil {
		logger.Errorf("Error sending summary of chat messages: %s", err.Error())
		return
	}

	logger.Infof("Jarvis Summarize responded with message: %s", resp.Choices[0].Message.Content)
}
