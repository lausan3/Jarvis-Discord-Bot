package help

import (
	"fmt"
	"main/infra/logger"

	"github.com/bwmarrin/discordgo"
)

func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	helpMessage := fmt.Sprintf("Certainly, %s! Please feel free to look at all the commands you can use below:\n\n", m.Author.Mention())

	helpMessage += CreateHelpCommandFormat("Jarvis hello", "responds to you with hello!")
	helpMessage += CreateHelpCommandFormat("Jarvis help", "this message.")
	helpMessage += CreateHelpCommandFormat("Jarvis summarize", "summarizes the last x messages. This amount can be changed in the options menu")
	helpMessage += CreateHelpCommandFormat("Jarvis summarize before Message ID", "summarizes the last x messages before a given message ID. This amount can be changed in the options menu.")
	// helpMessage += CreateHelpCommandFormat("Jarvis options", "opens the options menu.")

	if _, err := s.ChannelMessageSendReply(m.ChannelID, helpMessage, m.Reference()); err != nil {
		s.ChannelMessageSend(m.ChannelID, "Sorry, an error has occurred. Please try again.")
		logger.Errorf("Failed to send a help message for user %s. Message ID: %s", m.Author.Username, m.ID)
		return
	}

	guild, err := s.Guild(m.GuildID)

	if err != nil {
		logger.Errorf("Could not retrieve guild info for guild %s.", m.GuildID)
	}

	channel, err := s.Channel(m.ChannelID)

	if err != nil {
		logger.Errorf("Could not retrieve channel info for channel %s.", m.ChannelID)
	}

	logger.Infof(fmt.Sprintf("Successfully sent help message in channel %s.%s.", guild.Name, channel.Name))
}

func CreateHelpCommandFormat(commandName string, description string) string {
	return fmt.Sprintf("``%s`` - %s\n", commandName, description)
}
