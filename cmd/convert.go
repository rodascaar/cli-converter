package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/example/audioconv/internal/config"
	"github.com/example/audioconv/internal/converter"
	"github.com/example/audioconv/internal/ui"
	"github.com/spf13/cobra"
)

var (
	convertBitrate    string
	convertQuality    int
	convertSampleRate int
	convertChannels   string
	convertPreset     string
	convertNormalize  bool
	convertTrimStart  string
	convertTrimEnd    string
	convertOverwrite  bool
	convertDryRun     bool
	convertQuiet      bool
	convertVerbose    bool
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert <input> <output>",
	Short: "Convertir un solo archivo de audio",
	Long: `Convertir un archivo de audio a otro formato con control fino sobre
parámetros de calidad como bitrate, sample rate, canales, etc.

Ejemplos:
  audioconv convert song.flac song.mp3
  audioconv convert song.wav song.mp3 -b 320k -q 0
  audioconv convert podcast.m4a podcast.mp3 -c mono -b 64k
  audioconv convert music.flac music.opus -p high --normalize`,
	Args: cobra.ExactArgs(2),
	RunE: runConvert,
}

func init() {
	rootCmd.AddCommand(convertCmd)

	// Quality options
	convertCmd.Flags().StringVarP(&convertBitrate, "bitrate", "b", "", "bitrate de salida (ej: 192k, 320k)")
	convertCmd.Flags().IntVarP(&convertQuality, "quality", "q", 0, "calidad VBR 0-9 para mp3/ogg (0=mejor)")
	convertCmd.Flags().IntVarP(&convertSampleRate, "sample-rate", "s", 0, "frecuencia de muestreo (ej: 44100, 48000)")
	convertCmd.Flags().StringVarP(&convertChannels, "channels", "c", "auto", "canales: mono|stereo|auto")
	convertCmd.Flags().StringVarP(&convertPreset, "preset", "p", "", "preset de calidad: low|medium|high|ultra|lossless")

	// Processing options
	convertCmd.Flags().BoolVar(&convertNormalize, "normalize", false, "normalizar volumen a -16 LUFS")
	convertCmd.Flags().StringVar(&convertTrimStart, "trim-start", "", "recortar inicio (ej: 10s, 1m30s)")
	convertCmd.Flags().StringVar(&convertTrimEnd, "trim-end", "", "recortar desde el final")

	// Control options
	convertCmd.Flags().BoolVarP(&convertOverwrite, "overwrite", "y", false, "sobrescribir sin preguntar")
	convertCmd.Flags().BoolVar(&convertDryRun, "dry-run", false, "mostrar comando ffmpeg sin ejecutar")
	convertCmd.Flags().BoolVar(&convertQuiet, "quiet", false, "suprimir output excepto errores")
	convertCmd.Flags().BoolVarP(&convertVerbose, "verbose", "v", false, "mostrar salida completa de ffmpeg")
}

func runConvert(cmd *cobra.Command, args []string) error {
	input := args[0]
	output := args[1]

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		ui.PrintError(fmt.Sprintf("Config loading error: %v", err))
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Apply defaults from config
	opts := converter.ConvertOptions{
		OutputFormat: getOutputFormat(output),
		Bitrate:      convertBitrate,
		Quality:      convertQuality,
		SampleRate:   convertSampleRate,
		Channels:     convertChannels,
		Preset:       convertPreset,
		Normalize:    convertNormalize,
		TrimStart:    convertTrimStart,
		TrimEnd:      convertTrimEnd,
		Overwrite:    convertOverwrite,
		DryRun:       convertDryRun,
		Quiet:        convertQuiet,
		Verbose:      convertVerbose,
	}

	// Apply config defaults if not specified
	if opts.Bitrate == "" && cfg.Defaults.Bitrate != "" {
		opts.Bitrate = cfg.Defaults.Bitrate
	}
	if opts.Quality == 0 && cfg.Defaults.Quality > 0 {
		opts.Quality = cfg.Defaults.Quality
	}
	if opts.SampleRate == 0 && cfg.Defaults.SampleRate > 0 {
		opts.SampleRate = cfg.Defaults.SampleRate
	}
	if opts.Channels == "auto" && cfg.Defaults.Channels != "" {
		opts.Channels = cfg.Defaults.Channels
	}
	if !opts.Overwrite && cfg.Defaults.Overwrite {
		opts.Overwrite = cfg.Defaults.Overwrite
	}

	// Create converter
	conv := converter.NewConverter(cfg)

	// Check ffmpeg availability
	if _, err := converter.CheckFFmpegVersion(); err != nil {
		return fmt.Errorf("ffmpeg check failed: %w", err)
	}

	// Setup progress reporting
	var reporter ui.ProgressReporter
	if !opts.Quiet && !opts.DryRun {
		reporter = ui.NewProgressBar(fmt.Sprintf("Converting %s", filepath.Base(input)), 100, os.Stderr)
		reporter.Start(100)
	}

	// Perform conversion
	ctx := context.Background()
	start := time.Now()

	err = conv.Convert(ctx, input, output, opts)

	duration := time.Since(start)

	if reporter != nil {
		if err != nil {
			reporter.Error(err)
		} else {
			reporter.Finish()
		}
	}

	// Print result
	if !opts.Quiet {
		if err != nil {
			ui.PrintError(fmt.Sprintf("Conversion failed: %v", err))
			return err
		}

		ui.PrintSuccess(fmt.Sprintf("Converted %s → %s in %s",
			filepath.Base(input),
			filepath.Base(output),
			duration.Round(time.Second)))

		// Show file sizes if conversion succeeded
		if !opts.DryRun {
			showConversionStats(input, output)
		}
	}

	return nil
}

// getOutputFormat extracts format from output filename
func getOutputFormat(output string) string {
	ext := filepath.Ext(output)
	if ext == "" {
		return ""
	}
	return ext[1:] // Remove the dot
}

// showConversionStats shows before/after file sizes
func showConversionStats(input, output string) {
	inputInfo, err := os.Stat(input)
	if err != nil {
		return
	}

	outputInfo, err := os.Stat(output)
	if err != nil {
		return
	}

	inputSize := inputInfo.Size()
	outputSize := outputInfo.Size()
	saved := inputSize - outputSize
	savedPercent := float64(saved) / float64(inputSize) * 100

	fmt.Printf("  Original: %s | New: %s | Saved: %s (%.1f%%)\n",
		ui.FormatSize(inputSize),
		ui.FormatSize(outputSize),
		ui.FormatSize(saved),
		savedPercent)
}
