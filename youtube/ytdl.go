package youtube

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"go.uber.org/zap"
)

const AudioFormat = "opus"
const DownloadDir = "downloads"
const YTDLPDownloadFormat = "%(id)s.%(ext)s"

var downloadPath string = path.Join(DownloadDir, YTDLPDownloadFormat)

func DownloadAudio(videoUrl string) (string, error) {
	_, videoTag, _ := strings.Cut(videoUrl, "?v=")
	videoTagContainsAmpersand := strings.Contains(videoTag, "&")
	if videoTagContainsAmpersand {
		videoTag, _, _ = strings.Cut(videoTag, "&")
	}

	audioPath := path.Join(DownloadDir, fmt.Sprintf("%s.%s", videoTag, AudioFormat))

	if _, err := os.Stat(audioPath); err == nil {
		zap.L().Sugar().Info(fmt.Sprintf("Videotag found in cache: %s", videoTag))
		return audioPath, nil
	}

	zap.L().Sugar().Info(fmt.Sprintf("Downloading videotag %s", videoTag))

	cmd := exec.Command("yt-dlp", "-x", videoUrl, "-o", downloadPath)
	_, err := cmd.Output()

	if err != nil {
		zap.L().Sugar().Error("Error when calling yt-dlp:", err)
		return "", err
	}

	return audioPath, nil
}
