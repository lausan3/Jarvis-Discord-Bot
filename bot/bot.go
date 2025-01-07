package bot

import (
	"fmt"
	"main/bot/messagecommands/help"
	"main/bot/messagecommands/summarize"
	"main/bot/messagecommands/testresponse"
	"main/infra/logger"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

func Start() {
	appToken := viper.GetString("DISCORD_TOKEN")
	d, err := discordgo.New("Bot " + appToken)
	if err != nil {
		logger.Errorf("Discord bot instance could not be started, %s", err.Error())
	}

	// Add Handlers
	d.AddHandler(ready)

	d.AddHandler(messageCreate)

	// Add Intents
	d.Identify.Intents = discordgo.IntentsGuildMessages

	err = d.Open()
	defer d.Close()

	if err != nil {
		logger.Fatalf("Error opening Discord session: %s", err.Error())
	}

	fmt.Println("Jarvis is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

func ready(s *discordgo.Session, r *discordgo.Ready) {
	logger.Infof("STARTED Jarvis at time %s", time.Now().Format(time.DateTime))
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// The bot created the message, return early.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Look for jarvis prefix
	if strings.HasPrefix(strings.ToLower(m.Content), "jarvis") {
		commandParams := strings.Split(m.Content, " ")[1:]
		commandParamsLength := len(commandParams)

		if commandParamsLength == 0 {
			logger.Infof("RECEIVED message without input from user %s", m.Author.GlobalName)
			s.ChannelMessageSendReply(m.ChannelID, "Did you mean to give me a command? Type jarvis help for a list of my available commands.", m.MessageReference)
			return
		}

		switch strings.ToLower(commandParams[0]) {
		case "hello", "test":
			testresponse.TestResponse(s, m)
		case "summarize":
			// Normal summarize command
			if commandParamsLength == 1 {
				summarize.SummarizeBeforeMessageID(s, m, m.ID)
				// Summarize before command
			} else if commandParamsLength >= 3 && commandParams[1] == "before" {
				summarize.SummarizeBeforeMessageID(s, m, commandParams[2])
			}
		case "help":
			help.Help(s, m)
		default:
			logger.Infof("RECEIVED unknown message at %s: %s", m.Timestamp, m.Content)
		}
	}
}
