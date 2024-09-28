package discordbot

import (
	"fmt"

	"github.com/Rmkek/goyt/youtube"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var stopPlaying chan bool
var voiceConn *discordgo.VoiceConnection
var playing bool = false

func init() {
	stopPlaying = make(chan bool, 1)
}

func playSound(soundFile string, s *discordgo.Session, guildID, channelID string, stop <-chan bool) (err error) {
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	voiceConn = vc
	if err != nil {
		return err
	}

	playing = true
	dgvoice.PlayAudioFile(vc, soundFile, stop)
	playing = false

	return nil
}

func PlayHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	youtubeLink := options[0].Value.(string)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Downloading link: %s", youtubeLink),
		},
	})

	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		zap.L().Sugar().Error("Couldn't find channel guild", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Couldn't join voice channel.",
		})
		return
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == i.Member.User.ID {
			audioPath, err := youtube.DownloadAudio(youtubeLink)
			if err != nil {
				zap.L().Sugar().Error("Couldn't download video", err)
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Couldn't download video.",
				})
				return
			}

			if !playing {
				go playSound(audioPath, s, g.ID, vs.ChannelID, stopPlaying)
			}
			// TODO: else send song to soundQueue
		}
	}
}

func StopHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	stopPlaying <- true
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Playing stopped.",
		},
	})
}

func QuitHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var quitResponse string

	if !voiceConn.Ready {
		quitResponse = "Not connected to voice"
	} else {
		quitResponse = "Bye-bye :("
		voiceConn.Disconnect()
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: quitResponse,
		},
	})
}
