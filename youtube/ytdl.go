package youtube

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"strings"
)

const AudioFormat = "opus"
const DownloadDir = "downloads"
const YTDLPDownloadFormat = "%(id)s.%(ext)s"

var downloadPath string = path.Join(DownloadDir, YTDLPDownloadFormat)

func DownloadAudio(videoUrl string) (string, error) {
	_, videoTag, _ := strings.Cut(videoUrl, "?v=")
	audioPath := path.Join(DownloadDir, fmt.Sprintf("%s.%s", videoTag, AudioFormat))

	if _, err := os.Stat(audioPath); errors.Is(err, os.ErrNotExist) {
		slog.Info("Downloading videotag", "videoTag", videoTag)

		cmd := exec.Command("yt-dlp", "-x", videoUrl, "-o", downloadPath)
		_, err := cmd.Output()

		if err != nil {
			slog.Error("Error when calling yt-dlp:", slog.Any("error", err))
			return "", err
		}
	} else {
		slog.Info("Videotag found in cache.", "videoTag", videoTag)
	}

	return audioPath, nil
}
