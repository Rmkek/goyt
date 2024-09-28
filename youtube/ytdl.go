package youtube

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"go.uber.org/zap"
)

const (
	AudioFormat         = "opus"
	DownloadDir         = "downloads"
	YTDLPDownloadFormat = "%(id)s.%(ext)s"
)

var downloadPath string = path.Join(DownloadDir, YTDLPDownloadFormat)

func parseVideoURL(videoURL string) string {
	_, videoTag, _ := strings.Cut(videoURL, "?v=")
	videoTagContainsAmpersand := strings.Contains(videoTag, "&")

	if videoTagContainsAmpersand {
		videoTag, _, _ = strings.Cut(videoTag, "&")
	}

	return videoTag
}

func DownloadAudio(videoURL string) (string, error) {
	videoTag := parseVideoURL(videoURL)

	audioPath := path.Join(DownloadDir, fmt.Sprintf("%s.%s", videoTag, AudioFormat))

	if _, err := os.Stat(audioPath); err == nil {
		zap.L().Sugar().Info(fmt.Sprintf("Videotag found in cache: %s", videoTag))
		return audioPath, nil
	}

	zap.L().Sugar().Info(fmt.Sprintf("Downloading videotag %s", videoTag))

	cmd := exec.Command("yt-dlp", "-x", videoURL, "-o", downloadPath)
	_, err := cmd.Output()
	if err != nil {
		zap.L().Sugar().Error("Error when calling yt-dlp:", err)
		return "", err
	}

	return audioPath, nil
}
