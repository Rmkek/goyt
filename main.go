package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/Rmkek/goyt/discordbot"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var bot *discordgo.Session

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "play",
			Description: "Plays YouTube video link that you provide",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "youtube-link",
					Description: "Youtube link",
					Required:    true,
				},
			},
		},
		{
			Name:        "stop",
			Description: "Stops the player",
		},
		{
			Name:        "quit",
			Description: "Quits the player from voice channel",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"play": discordbot.PlayHandler,
		"stop": discordbot.StopHandler,
		"quit": discordbot.QuitHandler,
	}
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	discordBotToken := os.Getenv("DISCORD_BOT_TOKEN")

	bot, err = discordgo.New(fmt.Sprintf("Bot %s", discordBotToken))
	if err != nil {
		log.Fatal("Bot session failed to start", err)
	}
	bot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	logger := zap.NewExample()
	defer logger.Sync()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	bot.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		zap.L().Sugar().Infof("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	if err := bot.Open(); err != nil {
		zap.L().Sugar().Fatalf("Cannot open the session: %v", err)
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := bot.ApplicationCommandCreate(bot.State.User.ID, "", v)
		if err != nil {
			zap.L().Sugar().Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		zap.L().Sugar().Infof("Registered '%v' command.", v.Name)
		registeredCommands[i] = cmd
	}

	defer bot.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	zap.L().Sugar().Info("Press Ctrl+C to exit")
	<-stop
}
