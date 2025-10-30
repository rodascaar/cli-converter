package ui

import (
	"github.com/fatih/color"
)

// Color functions for consistent output
var (
	Success  = color.New(color.FgGreen).SprintFunc()
	Error    = color.New(color.FgRed).SprintFunc()
	Warning  = color.New(color.FgYellow).SprintFunc()
	Info     = color.New(color.FgBlue).SprintFunc()
	Progress = color.New(color.FgCyan).SprintFunc()

	// Symbols
	SuccessSymbol = Success("✓")
	ErrorSymbol   = Error("✗")
	WarningSymbol = Warning("⚠")
	InfoSymbol    = Info("ℹ")
)

// PrintSuccess prints a success message
func PrintSuccess(message string) {
	color.Green("%s %s", SuccessSymbol, message)
}

// PrintError prints an error message
func PrintError(message string) {
	color.Red("%s %s", ErrorSymbol, message)
}

// PrintWarning prints a warning message
func PrintWarning(message string) {
	color.Yellow("%s %s", WarningSymbol, message)
}

// PrintInfo prints an info message
func PrintInfo(message string) {
	color.Blue("%s %s", InfoSymbol, message)
}
