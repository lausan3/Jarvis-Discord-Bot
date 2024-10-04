package messagecommands

import (
	"context"
	"fmt"
	"main/infra/logger"
	"slices"
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

	logger.Infof("Successfully sent response to test command to user %s", m.Author.GlobalName)
}

// The same as Summarize() except for messages before a certain message id.
func SummarizeBeforeMessageID(s *discordgo.Session, m *discordgo.MessageCreate, beforeMessageID string) {
	logger.Infof("Received summarize command from user %s", m.Author.GlobalName)

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

	messagesArr, err := s.ChannelMessages(m.ChannelID, 10, beforeMessageID, "", "")
	if err != nil {
		logger.Errorf("Error getting messages before message id %s: %s", m.ID, err.Error())
		s.ChannelMessageSend(m.ChannelID, "You didn't provide a message id correctly, try again?")
		return
	}

	messages := []string{}

	for _, message := range messagesArr {
		isMessageFromBot := message.Author.ID == s.State.User.ID
		messageIsCommand := strings.HasPrefix(strings.ToLower(message.Content), "jarvis")
		messageIsEmpty := message.Author.GlobalName == "" || message.Content == ""

		if isMessageFromBot || messageIsCommand || messageIsEmpty {
			continue
		}

		messages = append(messages, fmt.Sprintf("%s said: %s", message.Author.GlobalName, message.Content))
	}

	if len(messages) == 0 {
		logger.Warnf("No messages found before message %s", m.ID)
		s.ChannelMessageSend(m.ChannelID, "I couldn't find any messages to summarize before your provided message. Please note that I don't summarize any jarvis commands or any message I send.")
		return
	}

	slices.Reverse(messages)

	logger.Infof("SENDING %v to ChatGPT to summarize", messages)

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
