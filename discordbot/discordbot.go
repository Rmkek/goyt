package discordbot

import (
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/Rmkek/goyt/youtube"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

var RmkID string = "253899307332272128"

func Ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateGameStatus(0, "!goyt")
}

func parseSoundPlayRequest(msg string) (string, error) {
	after, found := strings.CutPrefix(msg, "!goyt")

	if !found {
		return "", errors.New("no !goyt found in command")
	}

	return strings.TrimSpace(after), nil
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

	if strings.HasPrefix(m.Content, "!goyt") {
		// Find the channel that the message came from.
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			slog.Error("Couldn't find channel where the message came from", slog.Any("error", err))
			return
		}

		// Find the guild for that channel.
		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			slog.Error("Couldn't find channel guild", slog.Any("error", err))
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

				videoUrl, err := parseSoundPlayRequest(m.Content)
				if err != nil {
					slog.Error("Error when parsing sound play request", slog.Any("error", err))
					return
				}

				audioPath, err := youtube.DownloadAudio(videoUrl)
				if err != nil {
					slog.Error("Couldn't download video", slog.Any("error", err))
				}

				err = PlaySound(audioPath, s, g.ID, vs.ChannelID)

				if err != nil {
					slog.Error("Couldn't play sound", slog.Any("error", err))
				}

				return
			}
		}
	}
}

func PlaySound(soundFile string, s *discordgo.Session, guildID, channelID string) (err error) {
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	time.Sleep(250 * time.Millisecond)
	vc.Speaking(true)
	dgvoice.PlayAudioFile(vc, soundFile, make(chan bool))
	time.Sleep(250 * time.Millisecond)
	vc.Speaking(false)
	vc.Disconnect()

	return nil
}
