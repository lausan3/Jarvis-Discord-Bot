package bot

import (
	"fmt"
	"main/bot/messagecommands"
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
	logger.Infof("Jarvis started at time %s", time.Now().Format(time.DateTime))
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	logger.Infof("Message received at %s: %s", m.Timestamp, m.Content)

	// The bot created the message, return early.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Look for jarvis prefix
	if strings.HasPrefix(m.Content, "jarvis") {
		contentWithoutJarvis := strings.Split(m.Content, " ")[1:]

		if len(contentWithoutJarvis) == 0 {
			logger.Errorf("Message received without input from user %s", m.Author.GlobalName)
			return
		}

		if contentWithoutJarvis[0] == "test" {
			messagecommands.Test(s, m)
		}
	}
}
