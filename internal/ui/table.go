package ui

import (
	"fmt"
	"time"
)

// PrintConversionTable prints a table of conversion results
func PrintConversionTable(results []ConversionResult) {
	if len(results) == 0 {
		return
	}

	// Simple text output instead of tablewriter for now
	fmt.Println("Conversion Results:")
	fmt.Println("==================")

	var totalOriginal, totalNew int64
	var successCount, failCount int

	for _, result := range results {
		status := Success("✓ Success")
		if result.Error != nil {
			status = Error("✗ Failed")
			failCount++
		} else {
			successCount++
			totalOriginal += result.OriginalSize
			totalNew += result.NewSize
		}

		originalSize := FormatSize(result.OriginalSize)
		newSize := FormatSize(result.NewSize)
		timeStr := result.Duration.Round(time.Second).String()

		fmt.Printf("%s | %s | %s | %s | %s\n",
			result.Filename,
			status,
			originalSize,
			newSize,
			timeStr)
	}

	// Print summary
	if len(results) > 1 {
		saved := totalOriginal - totalNew
		savedPercent := float64(saved) / float64(totalOriginal) * 100

		fmt.Println()
		PrintInfo(fmt.Sprintf("Total: %d files | Success: %d | Failed: %d",
			len(results), successCount, failCount))

		if successCount > 0 {
			PrintInfo(fmt.Sprintf("Space saved: %s (%.1f%%)",
				FormatSize(saved), savedPercent))
		}
	}
}

// ConversionResult represents the result of a conversion
type ConversionResult struct {
	Filename     string
	Error        error
	OriginalSize int64
	NewSize      int64
	Duration     time.Duration
}
