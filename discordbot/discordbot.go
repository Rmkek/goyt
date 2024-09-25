package discordbot

import (
	"errors"
	"strings"

	"github.com/Rmkek/goyt/youtube"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var RmkID string = "253899307332272128"

func Ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateGameStatus(0, "!goyt")
}

func parseVideoUrlFromRequest(msg string) (videoUrl string) {
	return strings.TrimSpace(msg)
}

func parseBotRequest(botRequest string) (string, error) {
	after, found := strings.CutPrefix(botRequest, "!goyt")

	if !found {
		return "", errors.New("!goyt missing from request")
	}

	return after, nil
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Author.ID != RmkID {
		return
	}

	botRequest, err := parseBotRequest(m.Content)
	if err != nil {
		zap.L().Sugar().Error("Parse bot request error", err)
		return
	}

	// Find the channel that the message came from.
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		zap.L().Sugar().Error("Couldn't find channel where the message came from", err)
		return
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		zap.L().Sugar().Error("Couldn't find channel guild", err)
		return
	}

	// Look for the message sender in that guild's current voice states.
	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			if strings.Contains(m.Content, "stop") {
				vc, err := s.ChannelVoiceJoin(c.GuildID, vs.ChannelID, false, true)
				if err != nil {
					return
				}

				vc.Speaking(false)
				vc.Disconnect()
				return
			}

			videoUrl := parseVideoUrlFromRequest(botRequest)
			audioPath, err := youtube.DownloadAudio(videoUrl)
			if err != nil {
				zap.L().Sugar().Error("Couldn't download video", err)
			}

			err = PlaySound(audioPath, s, g.ID, vs.ChannelID)

			if err != nil {
				zap.L().Sugar().Error("Couldn't play sound", err)
			}

			return
		}
	}
}

func PlaySound(soundFile string, s *discordgo.Session, guildID, channelID string) (err error) {
	zap.L().Sugar().Infoln("VCs:", s.VoiceConnections)
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}
	zap.L().Sugar().Infoln("VCs:", s.VoiceConnections)

	vc.Speaking(true)
	dgvoice.PlayAudioFile(vc, soundFile, make(chan bool))
	vc.Speaking(false)
	vc.Disconnect()

	return nil
}
