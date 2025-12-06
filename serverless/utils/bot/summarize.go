package bot_utils

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

func SummarizeBefore(botToken string, openAIToken string, beforeMessageID string, s *discordgo.InteractionCreate) (response string, err error) {
	c, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return
	}

	if openAIToken == "" {
		err = errors.New("No OpenAI token provided")
		return
	}

	oai := openai.NewClient(openAIToken)

	channelId := s.ChannelID
	messagesArr, err := c.ChannelMessages(channelId, 50, beforeMessageID, "", "")
	if err != nil {
		c.ChannelMessageSend(channelId, "You didn't provide a message id correctly, try again?")
		return
	}

	messages := []string{}

	for _, message := range messagesArr {
		isMessageFromBot := message.Author.ID == s.AppID
		messageIsCommand := strings.HasPrefix(strings.ToLower(message.Content), "jarvis")
		messageIsEmpty := message.Author.GlobalName == "" || message.Content == ""

		if isMessageFromBot || messageIsCommand || messageIsEmpty {
			continue
		}

		messages = append(messages, fmt.Sprintf("%s said: %s", message.Author.GlobalName, message.Content))
	}

	if len(messages) == 0 {
		err = errors.New("No messages to summarize for messages before id " + beforeMessageID)
		c.ChannelMessageSend(channelId, "I couldn't find any messages to summarize before your provided message ID. Please note that I don't summarize any jarvis commands or any message I send.")
		return
	}

	slices.Reverse(messages)

	logrus.Infof("SENDING %v to ChatGPT to summarize", messages)

	systemPrompt := `
		You are Jarvis, a helpful chat summarization assistant bot on Discord.

		You will be given context in the form of a string representing messages in the format <USER NAME> said: <MESSAGE CONTENT>\n.
		Ignore any messages sent by YOU, Jarvis.

		Your task is to take this context and summarize what the conversation was about.

		Your summary should be concise, clear, and capture the main points of the conversation.
	`

	resp, err := oai.CreateChatCompletion(
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

	if err != nil || len(resp.Choices) == 0 {
		logrus.Errorf("OpenAI couldn't generate a response. Error: %v", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
