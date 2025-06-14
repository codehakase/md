package theme

import (
	"os"
	"runtime"
	"strings"
)

// BackgroundType represents the detected terminal background type
type BackgroundType int

const (
	// BackgroundUnknown indicates the background type could not be determined
	BackgroundUnknown BackgroundType = iota
	// BackgroundLight indicates a light terminal background
	BackgroundLight
	// BackgroundDark indicates a dark terminal background
	BackgroundDark
)

// String returns a string representation of the background type
func (bt BackgroundType) String() string {
	switch bt {
	case BackgroundLight:
		return "light"
	case BackgroundDark:
		return "dark"
	default:
		return "unknown"
	}
}

// DetectTerminalBackground attempts to detect the terminal background preference
// Returns BackgroundDark as the safe default for v1 implementation
func DetectTerminalBackground() BackgroundType {
	if theme := os.Getenv("COLORFGBG"); theme != "" {
		// COLORFGBG format is typically "foreground;background"
		// Lower numbers (0-7) typically indicate darker colors
		// Higher numbers (8-15) typically indicate lighter colors
		parts := strings.Split(theme, ";")
		if len(parts) >= 2 {
			bg := parts[len(parts)-1]
			// Simple heuristic: if background is 0-7, it's likely dark
			if bg >= "0" && bg <= "7" {
				return BackgroundDark
			} else if bg >= "8" && bg <= "15" {
				return BackgroundLight
			}
		}
	}

	if isDarkThemeEnvironment() {
		return BackgroundDark
	}

	if isLightThemeEnvironment() {
		return BackgroundLight
	}

	// Platform-specific detection attempts
	switch runtime.GOOS {
	case "darwin":
		return detectMacOSBackground()
	case "linux":
		return detectLinuxBackground()
	case "windows":
		return detectWindowsBackground()
	}

	// Default to dark theme as specified for v1
	return BackgroundDark
}

func isDarkThemeEnvironment() bool {
	darkIndicators := []string{
		"DARK_MODE=1",
		"THEME=dark",
		"COLOR_SCHEME=dark",
	}

	for _, indicator := range darkIndicators {
		parts := strings.Split(indicator, "=")
		if len(parts) == 2 {
			if os.Getenv(parts[0]) == parts[1] {
				return true
			}
		}
	}

	term := strings.ToLower(os.Getenv("TERM"))
	termProgram := strings.ToLower(os.Getenv("TERM_PROGRAM"))
	
	// Some terminals that commonly default to dark themes
	darkTerminals := []string{
		"alacritty",
		"kitty",
		"wezterm",
	}

	for _, darkTerm := range darkTerminals {
		if strings.Contains(term, darkTerm) || strings.Contains(termProgram, darkTerm) {
			return true
		}
	}

	return false
}

func isLightThemeEnvironment() bool {
	lightIndicators := []string{
		"LIGHT_MODE=1",
		"THEME=light",
		"COLOR_SCHEME=light",
	}

	for _, indicator := range lightIndicators {
		parts := strings.Split(indicator, "=")
		if len(parts) == 2 {
			if os.Getenv(parts[0]) == parts[1] {
				return true
			}
		}
	}

	return false
}

func detectMacOSBackground() BackgroundType {
	// Check for macOS specific environment variables
	if os.Getenv("TERM_PROGRAM") == "Apple_Terminal" {
		// Apple Terminal detection could be enhanced with AppleScript
		// For now, default to dark
		return BackgroundDark
	}

	if os.Getenv("TERM_PROGRAM") == "iTerm.app" {
		// iTerm2 detection could be enhanced with iTerm2's APIs
		// For now, default to dark
		return BackgroundDark
	}

	return BackgroundDark
}

func detectLinuxBackground() BackgroundType {
	// Check for common Linux desktop environment variables
	desktop := strings.ToLower(os.Getenv("XDG_CURRENT_DESKTOP"))
	
	// GNOME desktop environment
	if strings.Contains(desktop, "gnome") {
		// Could check gsettings for theme preference
		// gsettings get org.gnome.desktop.interface gtk-theme
		return BackgroundDark
	}

	// KDE desktop environment
	if strings.Contains(desktop, "kde") {
		// Could check kde theme settings
		return BackgroundDark
	}

	return BackgroundDark
}

func detectWindowsBackground() BackgroundType {
	// Check Windows Terminal
	if os.Getenv("WT_SESSION") != "" {
		// Windows Terminal detected
		return BackgroundDark
	}

	// Check for PowerShell
	if strings.Contains(strings.ToLower(os.Getenv("PSModulePath")), "powershell") {
		return BackgroundDark
	}

	return BackgroundDark
}