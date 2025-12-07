package terminal

import (
	"os"

	"golang.org/x/term"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

var (
	// ColorEnabled can be toggled by --no-color flag
	ColorEnabled = true
)

// Colorize wraps text in ANSI color codes if enabled and output is a terminal
func Colorize(text string, color string) string {
	if !ColorEnabled || !IsTerminal() {
		return text
	}
	return color + text + ColorReset
}

func Success(msg string) string {
	return Colorize("[✓] "+msg, ColorGreen)
}

func Warning(msg string) string {
	return Colorize("[!] "+msg, ColorYellow)
}

func Error(msg string) string {
	return Colorize("[✗] "+msg, ColorRed)
}

func Info(msg string) string {
	return Colorize("[i] "+msg, ColorBlue)
}

// IsTerminal checks if stdout is a terminal
func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}
