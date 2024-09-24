package main

import (
	"fmt"
	"os/signal"
	"syscall"

	"os"

	"go.uber.org/zap"

	"github.com/Rmkek/goyt/discordbot"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	logger := zap.NewExample()
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	err := godotenv.Load()
	if err != nil {
		zap.L().Sugar().Fatal("Error loading .env file")
	}

	discordBotToken := os.Getenv("DISCORD_BOT_TOKEN")

	// Create a new Discord session using the provided bot token.
	bot, err := discordgo.New(fmt.Sprintf("Bot %s", discordBotToken))

	if err != nil {
		zap.L().Sugar().Fatal("Bot session failed to start", err)
	}

	defer bot.Close()

	bot.AddHandler(discordbot.Ready)
	bot.AddHandler(discordbot.MessageCreate)

	// We need information about guilds (which includes their channels),
	// messages and voice states.
	bot.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Open the websocket and begin listening.
	err = bot.Open()
	if err != nil {
		zap.L().Sugar().Fatal("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	logger.Info("GoYT is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
