package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/example/audioconv/internal/config"
	"github.com/example/audioconv/internal/converter"
	"github.com/example/audioconv/internal/ui"
	"github.com/spf13/cobra"
)

var (
	batchOutputDir     string
	batchParallel      int
	batchRecursive     bool
	batchKeepStructure bool
)

// batchCmd represents the batch command
var batchCmd = &cobra.Command{
	Use:   "batch <pattern>",
	Short: "Convertir múltiples archivos que coincidan con un patrón",
	Long: `Convertir múltiples archivos de audio que coincidan con un patrón glob.
Soporta conversión en paralelo y mantiene estructura de directorios opcionalmente.

Ejemplos:
  audioconv batch '*.flac' --to mp3 -b 320k
  audioconv batch 'music/**/*.wav' --to opus -p high -o ~/converted -r
  audioconv batch '*.m4a' --to mp3 -c mono -j 8`,
	Args: cobra.ExactArgs(1),
	RunE: runBatch,
}

func init() {
	rootCmd.AddCommand(batchCmd)

	// Required flags
	batchCmd.Flags().StringP("to", "", "", "formato de destino (requerido)")
	_ = batchCmd.MarkFlagRequired("to")

	// Output options
	batchCmd.Flags().StringVarP(&batchOutputDir, "output-dir", "o", "", "directorio de salida")
	batchCmd.Flags().BoolVarP(&batchRecursive, "recursive", "r", false, "buscar archivos recursivamente")
	batchCmd.Flags().BoolVar(&batchKeepStructure, "keep-structure", false, "mantener estructura de carpetas")

	// Performance options
	batchCmd.Flags().IntVarP(&batchParallel, "parallel", "j", 4, "número de conversiones simultáneas")

	// Inherit convert flags
	batchCmd.Flags().StringVarP(&convertBitrate, "bitrate", "b", "", "bitrate de salida (ej: 192k, 320k)")
	batchCmd.Flags().IntVarP(&convertQuality, "quality", "q", 0, "calidad VBR 0-9 para mp3/ogg")
	batchCmd.Flags().IntVarP(&convertSampleRate, "sample-rate", "s", 0, "frecuencia de muestreo")
	batchCmd.Flags().StringVarP(&convertChannels, "channels", "c", "auto", "canales: mono|stereo|auto")
	batchCmd.Flags().StringVarP(&convertPreset, "preset", "p", "", "preset de calidad")
	batchCmd.Flags().BoolVar(&convertNormalize, "normalize", false, "normalizar volumen")
	batchCmd.Flags().StringVar(&convertTrimStart, "trim-start", "", "recortar inicio")
	batchCmd.Flags().StringVar(&convertTrimEnd, "trim-end", "", "recortar desde el final")
	batchCmd.Flags().BoolVarP(&convertOverwrite, "overwrite", "y", false, "sobrescribir sin preguntar")
	batchCmd.Flags().BoolVar(&convertDryRun, "dry-run", false, "mostrar comandos sin ejecutar")
	batchCmd.Flags().BoolVar(&convertQuiet, "quiet", false, "suprimir output")
	batchCmd.Flags().BoolVarP(&convertVerbose, "verbose", "v", false, "mostrar salida completa")
}

func runBatch(cmd *cobra.Command, args []string) error {
	pattern := args[0]
	outputFormat, _ := cmd.Flags().GetString("to")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Determine output directory
	outputDir := batchOutputDir
	if outputDir == "" {
		outputDir = cfg.Defaults.OutputDir
		if outputDir == "" {
			outputDir = "./converted"
		}
	}

	// Find matching files
	files, err := findFiles(pattern, batchRecursive)
	if err != nil {
		return fmt.Errorf("failed to find files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no files found matching pattern: %s", pattern)
	}

	if !convertQuiet {
		fmt.Printf("Found %d files to convert\n", len(files))
		if !convertDryRun {
			fmt.Printf("Converting to %s format...\n", strings.ToUpper(outputFormat))
		}
	}

	// Setup conversion options
	opts := converter.ConvertOptions{
		OutputFormat: outputFormat,
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

	// Apply config defaults
	applyConfigDefaults(&opts, cfg)

	// Create converter
	conv := converter.NewConverter(cfg)

	// Check ffmpeg availability
	if _, err := converter.CheckFFmpegVersion(); err != nil {
		return fmt.Errorf("ffmpeg check failed: %w", err)
	}

	// Setup progress reporting for batch
	var reporter ui.ProgressReporter
	if !convertQuiet && !convertDryRun {
		reporter = ui.NewProgressBar(fmt.Sprintf("Converting %d files", len(files)), int64(len(files)), os.Stderr)
		reporter.Start(int64(len(files)))
	}

	// Perform batch conversion
	ctx := context.Background()
	start := time.Now()

	results := conv.ConvertBatch(ctx, files, outputDir, opts)

	duration := time.Since(start)

	if reporter != nil {
		reporter.Finish()
	}

	// Print results
	if !convertQuiet {
		ui.PrintConversionTable(convertResults(results))
		fmt.Printf("\nTotal time: %s\n", duration.Round(time.Second))
	}

	return nil
}

// findFiles finds files matching a pattern
func findFiles(pattern string, _ bool) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	// Filter out directories
	var files []string
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}
		if !info.IsDir() {
			files = append(files, match)
		}
	}

	return files, nil
}

// applyConfigDefaults applies configuration defaults to options
func applyConfigDefaults(opts *converter.ConvertOptions, cfg *config.Config) {
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
	if batchParallel == 4 && cfg.Defaults.ParallelJobs > 0 {
		batchParallel = cfg.Defaults.ParallelJobs
	}
}

// convertResults converts converter results to UI results
func convertResults(results []converter.Result) []ui.ConversionResult {
	uiResults := make([]ui.ConversionResult, len(results))

	for i, result := range results {
		var originalSize, newSize int64

		if result.InputFile != "" {
			if info, err := os.Stat(result.InputFile); err == nil {
				originalSize = info.Size()
			}
		}

		if result.OutputFile != "" {
			if info, err := os.Stat(result.OutputFile); err == nil {
				newSize = info.Size()
			}
		}

		uiResults[i] = ui.ConversionResult{
			Filename:     filepath.Base(result.InputFile),
			Error:        result.Error,
			OriginalSize: originalSize,
			NewSize:      newSize,
			Duration:     result.Duration,
		}
	}

	return uiResults
}
