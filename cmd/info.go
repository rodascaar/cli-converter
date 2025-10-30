package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/example/audioconv/internal/audio"
	"github.com/example/audioconv/internal/converter"
	"github.com/spf13/cobra"
)

var (
	infoJson   bool
	infoFormat string
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info <file>",
	Short: "Mostrar información detallada del archivo de audio",
	Long: `Mostrar información detallada del archivo de audio incluyendo
codec, bitrate, duración, sample rate, canales y metadatos.

Ejemplos:
  audioconv info song.mp3
  audioconv info music.flac --json
  audioconv info podcast.m4a -f duration`,
	Args: cobra.ExactArgs(1),
	RunE: runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)

	infoCmd.Flags().BoolVar(&infoJson, "json", false, "salida en formato JSON")
	infoCmd.Flags().StringVarP(&infoFormat, "format", "f", "", "mostrar solo un campo (duration, bitrate, codec)")
}

func runInfo(cmd *cobra.Command, args []string) error {
	filename := args[0]

	// Validate input file
	if err := audio.ValidateInputFile(filename); err != nil {
		return fmt.Errorf("invalid input file: %w", err)
	}

	// Create converter and get info
	conv := converter.NewConverter(nil)
	info, err := conv.GetInfo(filename)
	if err != nil {
		return fmt.Errorf("failed to get audio info: %w", err)
	}

	// Output based on flags
	if infoJson {
		return outputJson(info)
	}

	if infoFormat != "" {
		return outputField(info, infoFormat)
	}

	// Default formatted output
	info.PrintAudioInfo()

	return nil
}

// outputJson outputs information in JSON format
func outputJson(info *audio.AudioInfo) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(info)
}

// outputField outputs only a specific field
func outputField(info *audio.AudioInfo, field string) error {
	switch field {
	case "duration":
		fmt.Println(info.Duration)
	case "bitrate":
		fmt.Println(info.Bitrate)
	case "codec":
		fmt.Println(info.Codec)
	case "format":
		fmt.Println(info.Format)
	case "sample_rate":
		fmt.Printf("%d\n", info.SampleRate)
	case "channels":
		fmt.Printf("%d\n", info.Channels)
	case "size":
		fmt.Printf("%d\n", info.Size)
	default:
		return fmt.Errorf("unknown field: %s", field)
	}
	return nil
}
