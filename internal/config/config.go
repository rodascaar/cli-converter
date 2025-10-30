package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Defaults struct {
		OutputFormat string `mapstructure:"output_format" yaml:"output_format"`
		Bitrate      string `mapstructure:"bitrate" yaml:"bitrate"`
		Quality      int    `mapstructure:"quality" yaml:"quality"`
		SampleRate   int    `mapstructure:"sample_rate" yaml:"sample_rate"`
		Channels     string `mapstructure:"channels" yaml:"channels"`
		Normalize    bool   `mapstructure:"normalize" yaml:"normalize"`
		Overwrite    bool   `mapstructure:"overwrite" yaml:"overwrite"`
		OutputDir    string `mapstructure:"output_dir" yaml:"output_dir"`
		ParallelJobs int    `mapstructure:"parallel_jobs" yaml:"parallel_jobs"`
	} `mapstructure:"defaults" yaml:"defaults"`

	Presets map[string]Preset `mapstructure:"presets" yaml:"presets"`

	FFmpeg struct {
		Path      string   `mapstructure:"path" yaml:"path"`
		ExtraArgs []string `mapstructure:"extra_args" yaml:"extra_args"`
	} `mapstructure:"ffmpeg" yaml:"ffmpeg"`
}

// Preset represents a quality preset
type Preset struct {
	Format           string `mapstructure:"format" yaml:"format"`
	Bitrate          string `mapstructure:"bitrate" yaml:"bitrate"`
	Channels         string `mapstructure:"channels" yaml:"channels"`
	Normalize        bool   `mapstructure:"normalize" yaml:"normalize"`
	CompressionLevel int    `mapstructure:"compression_level" yaml:"compression_level"`
}

// LoadConfig loads configuration from file
func LoadConfig() (*Config, error) {
	viper.SetConfigName("audioconv")
	viper.SetConfigType("yaml")

	// Add config paths
	viper.AddConfigPath(".")
	if home, err := os.UserHomeDir(); err == nil {
		viper.AddConfigPath(home)
	}

	// Set defaults
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Log the config file path for debugging
			configFile := viper.ConfigFileUsed()
			if configFile == "" {
				configFile = "unknown"
			}
			return nil, fmt.Errorf("error reading config file %s: %w", configFile, err)
		}
		// Config file not found, use defaults
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

// setDefaults sets the default configuration values
func setDefaults() {
	viper.SetDefault("defaults.output_format", "mp3")
	viper.SetDefault("defaults.bitrate", "192k")
	viper.SetDefault("defaults.quality", 2)
	viper.SetDefault("defaults.sample_rate", 44100)
	viper.SetDefault("defaults.channels", "auto")
	viper.SetDefault("defaults.normalize", false)
	viper.SetDefault("defaults.overwrite", false)
	viper.SetDefault("defaults.output_dir", "./converted")
	viper.SetDefault("defaults.parallel_jobs", 4)

	viper.SetDefault("ffmpeg.path", "ffmpeg")
	viper.SetDefault("ffmpeg.extra_args", []string{})
}

// GetConfigFilePath returns the path to the config file
func GetConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".audioconv.yaml"), nil
}

// SaveDefaultConfig creates a default config file
func SaveDefaultConfig() error {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	defaultConfig := `defaults:
  output_format: mp3
  bitrate: 192k
  quality: 2
  sample_rate: 44100
  channels: auto
  normalize: false
  overwrite: false
  output_dir: ./converted
  parallel_jobs: 4

presets:
  podcast:
    format: mp3
    bitrate: 64k
    channels: mono
    normalize: true
  hifi:
    format: flac
    compression_level: 8
    sample_rate: 96000

ffmpeg:
  path: ffmpeg
  extra_args: []
`

	return os.WriteFile(configPath, []byte(defaultConfig), 0644)
}
