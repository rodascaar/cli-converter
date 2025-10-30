package converter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/example/audioconv/internal/audio"
	"github.com/example/audioconv/internal/config"
	"github.com/example/audioconv/internal/ui"
	"github.com/example/audioconv/pkg/formats"
)

// ConvertOptions represents conversion options
type ConvertOptions struct {
	OutputFormat string
	Bitrate      string
	Quality      int
	SampleRate   int
	Channels     string
	Preset       string
	Normalize    bool
	TrimStart    string
	TrimEnd      string
	Overwrite    bool
	DryRun       bool
	Quiet        bool
	Verbose      bool
}

// Result represents the result of a conversion
type Result struct {
	InputFile  string
	OutputFile string
	Error      error
	Duration   time.Duration
}

// Converter handles audio conversions
type Converter struct {
	config *config.Config
}

// NewConverter creates a new converter
func NewConverter(cfg *config.Config) *Converter {
	return &Converter{
		config: cfg,
	}
}

// Convert converts a single audio file
func (c *Converter) Convert(ctx context.Context, input, output string, opts ConvertOptions) error {
	// Validate input
	if err := audio.ValidateInputFile(input); err != nil {
		return fmt.Errorf("input validation failed: %w", err)
	}

	// Validate output
	if err := audio.ValidateOutputPath(output); err != nil {
		return fmt.Errorf("output validation failed: %w", err)
	}

	// Check if output exists and handle overwrite
	if !opts.Overwrite {
		if _, err := os.Stat(output); err == nil {
			return fmt.Errorf("output file exists: %s (use --overwrite to replace)", output)
		}
	}

	// Get input format
	inputFormat, _ := formats.GetFormatFromExtension(input)
	outputFormat := opts.OutputFormat
	if outputFormat == "" {
		outputFormat, _ = formats.GetFormatFromExtension(output)
	}

	// Validate format compatibility
	if err := audio.ValidateFormatCompatibility(inputFormat, outputFormat); err != nil {
		if !opts.Quiet {
			ui.PrintWarning(err.Error())
		}
	}

	// Build ffmpeg command
	args, err := c.buildFFmpegArgs(input, output, opts)
	if err != nil {
		return fmt.Errorf("failed to build ffmpeg command: %w", err)
	}

	if opts.DryRun {
		fmt.Printf("ffmpeg %s\n", strings.Join(args, " "))
		return nil
	}

	// Execute conversion with progress reporting
	return c.executeFFmpegWithProgress(ctx, args, opts)
}

// ConvertBatch converts multiple files
func (c *Converter) ConvertBatch(ctx context.Context, inputs []string, outputDir string, opts ConvertOptions) []Result {
	results := make([]Result, len(inputs))

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		for i := range results {
			results[i] = Result{Error: fmt.Errorf("failed to create output directory: %w", err)}
		}
		return results
	}

	// Process files sequentially for now (parallel implementation would go here)
	for i, input := range inputs {
		start := time.Now()

		// Generate output filename
		base := strings.TrimSuffix(filepath.Base(input), filepath.Ext(input))
		output := filepath.Join(outputDir, base+"."+opts.OutputFormat)

		err := c.Convert(ctx, input, output, opts)

		results[i] = Result{
			InputFile:  input,
			OutputFile: output,
			Error:      err,
			Duration:   time.Since(start),
		}
	}

	return results
}

// buildFFmpegArgs builds ffmpeg command arguments
func (c *Converter) buildFFmpegArgs(input, output string, opts ConvertOptions) ([]string, error) {
	args := []string{"-i", input}

	// Add format-specific options
	formatArgs, err := c.getFormatArgs(opts)
	if err != nil {
		return nil, err
	}
	args = append(args, formatArgs...)

	// Add metadata preservation
	args = append(args, "-map_metadata", "0")

	// Add output file
	args = append(args, "-y", output) // -y for overwrite

	return args, nil
}

// getFormatArgs returns format-specific ffmpeg arguments
func (c *Converter) getFormatArgs(opts ConvertOptions) ([]string, error) {
	var args []string

	format := opts.OutputFormat
	if format == "" {
		return args, nil
	}

	// Apply preset if specified
	if opts.Preset != "" {
		presetArgs, err := c.getPresetArgs(opts.Preset, format)
		if err != nil {
			return nil, err
		}
		args = append(args, presetArgs...)
	} else {
		// Manual quality settings
		if opts.Bitrate != "" {
			args = append(args, "-b:a", opts.Bitrate)
		}
		if opts.Quality > 0 {
			args = append(args, "-q:a", fmt.Sprintf("%d", opts.Quality))
		}
	}

	// Sample rate
	if opts.SampleRate > 0 {
		args = append(args, "-ar", fmt.Sprintf("%d", opts.SampleRate))
	}

	// Channels
	switch opts.Channels {
	case "mono":
		args = append(args, "-ac", "1")
	case "stereo":
		args = append(args, "-ac", "2")
	}

	// Normalization
	if opts.Normalize {
		args = append(args, "-af", "loudnorm=I=-16:TP=-1.5:LRA=11")
	}

	// Trimming
	if opts.TrimStart != "" {
		args = append(args, "-ss", opts.TrimStart)
	}
	if opts.TrimEnd != "" {
		args = append(args, "-to", opts.TrimEnd)
	}

	return args, nil
}

// getPresetArgs returns preset-specific arguments
func (c *Converter) getPresetArgs(preset, format string) ([]string, error) {
	switch preset {
	case "low":
		switch format {
		case "mp3":
			return []string{"-c:a", "libmp3lame", "-q:a", "6"}, nil
		case "opus":
			return []string{"-c:a", "libopus", "-b:a", "64k"}, nil
		case "aac":
			return []string{"-c:a", "aac", "-b:a", "128k"}, nil
		}
	case "medium":
		switch format {
		case "mp3":
			return []string{"-c:a", "libmp3lame", "-q:a", "4"}, nil
		case "opus":
			return []string{"-c:a", "libopus", "-b:a", "96k"}, nil
		case "aac":
			return []string{"-c:a", "aac", "-b:a", "192k"}, nil
		}
	case "high":
		switch format {
		case "mp3":
			return []string{"-c:a", "libmp3lame", "-q:a", "0"}, nil
		case "opus":
			return []string{"-c:a", "libopus", "-b:a", "128k", "-vbr", "on"}, nil
		case "aac":
			return []string{"-c:a", "aac", "-b:a", "256k"}, nil
		}
	case "lossless":
		switch format {
		case "flac":
			return []string{"-c:a", "flac", "-compression_level", "8"}, nil
		case "wav":
			return []string{"-c:a", "pcm_s16le"}, nil
		}
	}

	return nil, fmt.Errorf("unsupported preset '%s' for format '%s'", preset, format)
}

// executeFFmpegWithProgress executes ffmpeg with progress reporting
func (c *Converter) executeFFmpegWithProgress(ctx context.Context, args []string, _ ConvertOptions) error {
	return ExecuteFFmpegWithProgress(ctx, args, nil) // TODO: Pass proper progress reporter
}

// GetInfo gets audio file information
func (c *Converter) GetInfo(_ string) (*audio.AudioInfo, error) {
	return audio.GetAudioInfo("") // TODO: Implement properly
}
