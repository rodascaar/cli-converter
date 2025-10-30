package ui

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"

	"github.com/schollz/progressbar/v3"
)

// ProgressReporter interface for progress reporting
type ProgressReporter interface {
	Start(total int64)
	Update(current int64)
	Finish()
	Error(err error)
}

// ProgressBar wraps the progressbar library
type ProgressBar struct {
	bar    *progressbar.ProgressBar
	start  time.Time
	total  int64
	writer io.Writer
}

// NewProgressBar creates a new progress bar
func NewProgressBar(description string, total int64, writer io.Writer) *ProgressBar {
	if writer == nil {
		writer = io.Discard
	}

	bar := progressbar.NewOptions64(total,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(writer),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(50),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(writer, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
	)

	return &ProgressBar{
		bar:    bar,
		start:  time.Now(),
		total:  total,
		writer: writer,
	}
}

// Start starts the progress bar
func (p *ProgressBar) Start(total int64) {
	p.total = total
	p.start = time.Now()
	p.bar.ChangeMax64(total)
}

// Update updates the progress bar
func (p *ProgressBar) Update(current int64) {
	_ = p.bar.Set64(current)
}

// Finish finishes the progress bar
func (p *ProgressBar) Finish() {
	_ = p.bar.Finish()
	elapsed := time.Since(p.start)
	fmt.Fprintf(p.writer, "Completed in %s\n", elapsed.Round(time.Second))
}

// Error shows an error on the progress bar
func (p *ProgressBar) Error(err error) {
	p.bar.Describe(fmt.Sprintf("[ERROR] %v", err))
	_ = p.bar.Finish()
}

// ProgressParser parses ffmpeg progress output
type ProgressParser struct {
	reTime    *regexp.Regexp
	reSize    *regexp.Regexp
	reBitrate *regexp.Regexp
	total     time.Duration
	reporter  ProgressReporter
}

// NewProgressParser creates a new progress parser
func NewProgressParser(total time.Duration, reporter ProgressReporter) *ProgressParser {
	return &ProgressParser{
		reTime:    regexp.MustCompile(`time=(\d{2}):(\d{2}):(\d{2}\.\d{2})`),
		reSize:    regexp.MustCompile(`size=\s*(\d+)kB`),
		reBitrate: regexp.MustCompile(`bitrate=\s*(\d+\.\d+)kbits/s`),
		total:     total,
		reporter:  reporter,
	}
}

// ParseLine parses a line of ffmpeg output
func (p *ProgressParser) ParseLine(line string) {
	matches := p.reTime.FindStringSubmatch(line)
	if len(matches) == 4 {
		hours, _ := strconv.Atoi(matches[1])
		minutes, _ := strconv.Atoi(matches[2])
		seconds, _ := strconv.ParseFloat(matches[3], 64)

		current := time.Duration(hours)*time.Hour +
			time.Duration(minutes)*time.Minute +
			time.Duration(seconds)*time.Second

		if p.total > 0 {
			progress := int64(current) * 100 / int64(p.total)
			p.reporter.Update(progress)
		}
	}
}

// GetDurationFromFFmpegOutput extracts duration from ffmpeg output
func GetDurationFromFFmpegOutput(output string) (time.Duration, error) {
	re := regexp.MustCompile(`Duration: (\d{2}):(\d{2}):(\d{2}\.\d{2})`)
	matches := re.FindStringSubmatch(output)
	if len(matches) != 4 {
		return 0, fmt.Errorf("duration not found in output")
	}

	hours, _ := strconv.Atoi(matches[1])
	minutes, _ := strconv.Atoi(matches[2])
	seconds, _ := strconv.ParseFloat(matches[3], 64)

	return time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second, nil
}

// FormatDuration formats a duration for display
func FormatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	if h > 0 {
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	}
	return fmt.Sprintf("%dm%ds", m, s)
}

// FormatSize formats a size in bytes for display
func FormatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
