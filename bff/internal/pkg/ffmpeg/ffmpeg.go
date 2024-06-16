package ffmpeg

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
)

var ErrWrongTimeOutput = errors.New("invalid time output")

type FfmpegExecutor struct {
	log *zerolog.Logger
}

func New(log *zerolog.Logger) *FfmpegExecutor {
	return &FfmpegExecutor{
		log: log,
	}
}

type timespan time.Duration

func (f *FfmpegExecutor) GetAudioFromVideo(filename string) (string, error) {
	audioName := xid.New().String() + ".wav"
	flags := []string{
		"-i", filename,
		"-vn",
		"-acodec", "pcm_s16le",
		"-ar", "44100",
		"-ac", "2", audioName,
	}

	cmd := exec.Command("ffmpeg", flags...)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to run ffmpeg: %w", err)
	}

	return audioName, nil
}

func (t timespan) Format(format string) string {
	z := time.Unix(0, 0).UTC()
	return z.Add(time.Duration(t)).Format(format)
}

func (f *FfmpegExecutor) GetScreenshotFromVideo(filename string) (string, error) {
	id := xid.New().String() + ".png"

	length, err := f.GetVideoLength(filename)
	if err != nil {
		return "", fmt.Errorf("get video length failed: %w", err)
	}

	length /= 2

	flags := []string{"-ss", timespan(length).Format("15:04:05"), "-i", filename, "-frames:v", "1", id}

	cmd := exec.Command("ffmpeg", flags...)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg run failed: %w", err)
	}

	return id, nil
}

func (f *FfmpegExecutor) GetVideoLength(filename string) (time.Duration, error) {
	flags := []string{"-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", "-sexagesimal", filename}
	f.log.Debug().Strs("flags", flags).Msg("starting ffprobe")

	cmd := exec.Command("ffprobe", flags...)

	outputBytes, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe get output failed: %w", err)
	}

	return parseTime(string(outputBytes))
}

func parseTime(t string) (time.Duration, error) {
	sp := strings.Split(t, ".")
	if len(sp) != 2 {
		return 0, fmt.Errorf("%w: %s", ErrWrongTimeOutput, t)
	}

	tt, err := time.Parse("15:04:05", sp[0])
	if err != nil {
		return 0, fmt.Errorf("failed to parse time: %w", err)
	}

	return tt.Sub(time.Time{}), nil
}
