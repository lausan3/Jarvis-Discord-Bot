package messagecommands

import (
	"main/infra/logger"

	"github.com/bwmarrin/discordgo"
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
