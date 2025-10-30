package audio

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// AudioInfo represents audio file metadata
type AudioInfo struct {
	Format     string            `json:"format"`
	Codec      string            `json:"codec"`
	Duration   time.Duration     `json:"duration"`
	Bitrate    string            `json:"bitrate"`
	SampleRate int               `json:"sample_rate"`
	Channels   int               `json:"channels"`
	Size       int64             `json:"size"`
	Metadata   map[string]string `json:"metadata"`
}

// GetAudioInfo extracts metadata from an audio file using ffprobe
func GetAudioInfo(filename string) (*AudioInfo, error) {
	// Run ffprobe to get detailed info
	cmd := exec.Command("ffprobe",
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		filename)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run ffprobe: %w", err)
	}

	// Parse the JSON output (simplified parsing for now)
	info := &AudioInfo{
		Metadata: make(map[string]string),
	}

	// Extract basic info from ffprobe output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, `"codec_name"`) {
			// Extract codec
			re := regexp.MustCompile(`"codec_name":\s*"([^"]+)"`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				info.Codec = matches[1]
			}
		}
		if strings.Contains(line, `"duration"`) {
			re := regexp.MustCompile(`"duration":\s*"([^"]+)"`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if dur, err := strconv.ParseFloat(matches[1], 64); err == nil {
					info.Duration = time.Duration(dur * float64(time.Second))
				}
			}
		}
		if strings.Contains(line, `"bit_rate"`) {
			re := regexp.MustCompile(`"bit_rate":\s*"([^"]+)"`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if bitrate, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
					info.Bitrate = fmt.Sprintf("%d kbps", bitrate/1000)
				}
			}
		}
		if strings.Contains(line, `"sample_rate"`) {
			re := regexp.MustCompile(`"sample_rate":\s*"([^"]+)"`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if rate, err := strconv.Atoi(matches[1]); err == nil {
					info.SampleRate = rate
				}
			}
		}
		if strings.Contains(line, `"channels"`) {
			re := regexp.MustCompile(`"channels":\s*(\d+)`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if ch, err := strconv.Atoi(matches[1]); err == nil {
					info.Channels = ch
				}
			}
		}
		if strings.Contains(line, `"size"`) {
			re := regexp.MustCompile(`"size":\s*"([^"]+)"`)
			if matches := re.FindStringSubmatch(line); len(matches) > 1 {
				if size, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
					info.Size = size
				}
			}
		}
	}

	// Get file size if not found in ffprobe output
	if info.Size == 0 {
		if stat, err := exec.Command("stat", "-f%z", filename).Output(); err == nil {
			if size, err := strconv.ParseInt(strings.TrimSpace(string(stat)), 10, 64); err == nil {
				info.Size = size
			}
		}
	}

	// Determine format from codec
	info.Format = getFormatFromCodec(info.Codec)

	return info, nil
}

// getFormatFromCodec maps codec to format name
func getFormatFromCodec(codec string) string {
	switch codec {
	case "mp3":
		return "MP3"
	case "aac":
		return "AAC"
	case "flac":
		return "FLAC"
	case "pcm_s16le", "pcm_s24le", "pcm_s32le":
		return "WAV"
	case "vorbis":
		return "OGG Vorbis"
	case "opus":
		return "Opus"
	case "wmav2":
		return "WMA"
	case "alac":
		return "ALAC"
	default:
		return codec
	}
}

// PrintAudioInfo prints formatted audio information
func (info *AudioInfo) PrintAudioInfo() {
	fmt.Printf("Format:      %s (%s)\n", info.Format, info.Codec)
	fmt.Printf("Duration:    %s\n", formatDuration(info.Duration))
	fmt.Printf("Bitrate:     %s\n", info.Bitrate)
	fmt.Printf("Sample Rate: %d Hz\n", info.SampleRate)
	fmt.Printf("Channels:    %d (%s)\n", info.Channels, getChannelDescription(info.Channels))
	fmt.Printf("Size:        %s\n", formatSize(info.Size))

	if len(info.Metadata) > 0 {
		fmt.Println("\nMetadata:")
		for key, value := range info.Metadata {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}
}

// formatDuration formats duration for display
func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60

	if h > 0 {
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	}
	return fmt.Sprintf("%dm%ds", m, s)
}

// formatSize formats size for display
func formatSize(bytes int64) string {
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

// getChannelDescription returns human-readable channel description
func getChannelDescription(channels int) string {
	switch channels {
	case 1:
		return "Mono"
	case 2:
		return "Stereo"
	default:
		return fmt.Sprintf("%d channels", channels)
	}
}
