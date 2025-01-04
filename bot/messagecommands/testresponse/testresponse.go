package testresponse

import (
	"fmt"
	"main/infra/logger"

	"github.com/bwmarrin/discordgo"
)

// Send a test message as a response to "jarvis test"
func TestResponse(s *discordgo.Session, m *discordgo.MessageCreate) {
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Hello, %s!", m.Author.Mention()))
	if err != nil {
		logger.Errorf("Could not send test message response!")
		return
	}

	logger.Infof("Successfully sent response to test command to user %s", m.Author.GlobalName)
}
