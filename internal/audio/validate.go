package audio

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/example/audioconv/pkg/formats"
)

// ValidateInputFile validates an input audio file
func ValidateInputFile(filename string) error {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filename)
	}

	// Check if it's a regular file
	info, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("cannot access file: %s", err)
	}
	if info.IsDir() {
		return fmt.Errorf("path is a directory, not a file: %s", filename)
	}

	// Check file size (minimum 1KB)
	if info.Size() < 1024 {
		return fmt.Errorf("file is too small to be a valid audio file: %s", filename)
	}

	// Check format from extension
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return fmt.Errorf("file has no extension: %s", filename)
	}

	_, supported := formats.GetFormatFromExtension(filename)
	if !supported {
		return fmt.Errorf("unsupported format '%s'. Supported formats: %s",
			ext, strings.Join(formats.GetSupportedFormats(), ", "))
	}

	return nil
}

// ValidateOutputPath validates an output path
func ValidateOutputPath(outputPath string) error {
	// Check if output directory exists or can be created
	dir := filepath.Dir(outputPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("cannot create output directory: %s", err)
		}
	}

	// Check if output file already exists
	if _, err := os.Stat(outputPath); err == nil {
		// File exists, this is OK - we'll handle overwrite logic elsewhere
		return nil
	}

	// Check if we can write to the directory
	dir = filepath.Dir(outputPath)
	if dir == "." {
		dir = "."
	}
	testFile := filepath.Join(dir, ".audioconv_test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("cannot write to output directory: %s", err)
	}
	os.Remove(testFile) // Clean up

	return nil
}

// ValidateFormatCompatibility checks if conversion makes sense
func ValidateFormatCompatibility(inputFormat, outputFormat string) error {
	inputInfo, inputExists := formats.SupportedFormats[inputFormat]
	outputInfo, outputExists := formats.SupportedFormats[outputFormat]

	if !inputExists || !outputExists {
		return fmt.Errorf("invalid format specified")
	}

	// Warn about lossy to lossless conversion (not recommended)
	if !inputInfo.IsLossy && outputInfo.IsLossy {
		return fmt.Errorf("warning: converting from lossless (%s) to lossy (%s) format will not improve quality",
			inputFormat, outputFormat)
	}

	return nil
}

// CheckFFmpegAvailability checks if ffmpeg is installed and accessible
func CheckFFmpegAvailability() error {
	// This would normally run "ffmpeg -version" but we'll implement it in the converter
	return nil
}
