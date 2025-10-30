package converter

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/example/audioconv/internal/ui"
)

// FFmpegWrapper wraps ffmpeg execution with progress reporting
type FFmpegWrapper struct {
	// Fields removed as they were unused
}

// ExecuteFFmpegWithProgress executes ffmpeg with progress reporting
func ExecuteFFmpegWithProgress(ctx context.Context, args []string, reporter ui.ProgressReporter) error {
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)

	// Get pipes for stdout and stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	// Parse progress from stderr
	go parseProgress(stderr, reporter)

	// Wait for completion
	if err := cmd.Wait(); err != nil {
		reporter.Error(err)
		return fmt.Errorf("ffmpeg execution failed: %w", err)
	}

	reporter.Finish()
	return nil
}

// parseProgress parses ffmpeg progress output
func parseProgress(stderr io.ReadCloser, reporter ui.ProgressReporter) {
	scanner := bufio.NewScanner(stderr)
	parser := ui.NewProgressParser(0, reporter) // Duration will be set later

	for scanner.Scan() {
		line := scanner.Text()

		// Try to extract duration first
		if strings.Contains(line, "Duration:") {
			if duration, err := ui.GetDurationFromFFmpegOutput(line); err == nil {
				parser = ui.NewProgressParser(duration, reporter)
			}
		}

		// Parse progress
		parser.ParseLine(line)
	}
}

// CheckFFmpegVersion checks if ffmpeg is available and gets version
func CheckFFmpegVersion() (string, error) {
	cmd := exec.Command("ffmpeg", "-version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ffmpeg not found in PATH: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		// First line contains version info
		return strings.TrimSpace(lines[0]), nil
	}

	return "unknown", nil
}

// ValidateFFmpegCapabilities validates that ffmpeg has required codecs
func ValidateFFmpegCapabilities() error {
	cmd := exec.Command("ffmpeg", "-codecs")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("cannot check ffmpeg codecs: %w", err)
	}

	requiredCodecs := []string{
		"libmp3lame",
		"libfdk_aac",
		"flac",
		"libvorbis",
		"libopus",
	}

	outputStr := string(output)
	for _, codec := range requiredCodecs {
		if !strings.Contains(outputStr, codec) {
			return fmt.Errorf("required codec '%s' not available in ffmpeg", codec)
		}
	}

	return nil
}
