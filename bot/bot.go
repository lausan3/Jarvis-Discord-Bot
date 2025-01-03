package bot

import (
	"fmt"
	testresponse "main/bot/messagecommands/TestResponse"
	"main/bot/messagecommands/summarize"
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

	d.AddHandler(ready)

	d.AddHandler(messageCreate)

	d.Identify.Intents = discordgo.IntentsGuildMessages

	err = d.Open()
	if err != nil {
		logger.Fatalf("Error opening Discord session: %s", err.Error())
	}

	fmt.Println("Jarvis is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	defer d.Close()
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
			logger.Errorf("RECEIVED message without input from user %s", m.Author.GlobalName)
			s.ChannelMessageSendReply(m.ChannelID, "Did you mean to give me a command? Type jarvis help for a list of my available commands.", m.MessageReference)
			return
		}

		switch strings.ToLower(commandParams[0]) {
		case "test":
			testresponse.TestResponse(s, m)
		case "summarize":
			// Normal summarize command
			if commandParamsLength == 1 {
				summarize.SummarizeBeforeMessageID(s, m, m.ID)
				// Summarize before command
			} else if commandParamsLength >= 3 && commandParams[1] == "before" {
				summarize.SummarizeBeforeMessageID(s, m, commandParams[2])
			}
		default:
			logger.Infof("RECEIVED unknown message at %s: %s", m.Timestamp, m.Content)
		}
	}
}
