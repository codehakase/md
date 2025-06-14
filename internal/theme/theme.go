// Package theme provides terminal background detection and theme management
// for the md CLI tool. It automatically detects whether the terminal has
// a light or dark background and applies appropriate color schemes for
// both markdown rendering and code syntax highlighting.
//
// The package integrates with Chroma for syntax highlighting and provides
// consistent ANSI color styling across all markdown components.

package theme

import (
	"fmt"
	"os"
	"strings"
)

// ThemeManager manages styling and theming for the markdown renderer
type ThemeManager struct {
	backgroundType BackgroundType
	colors         map[string]string
	chromaTheme    string
}

// TerminalTheme represents terminal-specific theme information
type TerminalTheme struct {
	Background     BackgroundType
	ChromaTheme    string
	ColorScheme    string
	IsHighContrast bool
}

// ColorKey represents different styling elements
type ColorKey string

const (
	// Header colors
	Header1 ColorKey = "header1"
	Header2 ColorKey = "header2"
	Header3 ColorKey = "header3"
	Header4 ColorKey = "header4"
	Header5 ColorKey = "header5"
	Header6 ColorKey = "header6"

	// Text styling
	Bold          ColorKey = "bold"
	Italic        ColorKey = "italic"
	Strikethrough ColorKey = "strikethrough"
	Code          ColorKey = "code"

	// Block elements
	BlockQuote ColorKey = "blockquote"
	Link       ColorKey = "link"

	// List elements
	BulletPoint ColorKey = "bullet"
	OrderedList ColorKey = "ordered"

	// Table elements
	TableHeader ColorKey = "table_header"
	TableBorder ColorKey = "table_border"

	// Special
	Reset ColorKey = "reset"
)

// ANSI color constants for direct use
const (
	ANSIReset         = "\033[0m"
	ANSIBold          = "\033[1m"
	ANSIDim           = "\033[2m"
	ANSIItalic        = "\033[3m"
	ANSIUnderline     = "\033[4m"
	ANSIStrikethrough = "\033[9m"

	// Basic colors
	ANSIBlack   = "\033[30m"
	ANSIRed     = "\033[31m"
	ANSIGreen   = "\033[32m"
	ANSIYellow  = "\033[33m"
	ANSIBlue    = "\033[34m"
	ANSIMagenta = "\033[35m"
	ANSICyan    = "\033[36m"
	ANSIWhite   = "\033[37m"

	// Bright colors
	ANSIBrightBlack   = "\033[90m"
	ANSIBrightRed     = "\033[91m"
	ANSIBrightGreen   = "\033[92m"
	ANSIBrightYellow  = "\033[93m"
	ANSIBrightBlue    = "\033[94m"
	ANSIBrightMagenta = "\033[95m"
	ANSIBrightCyan    = "\033[96m"
	ANSIBrightWhite   = "\033[97m"
)

// New creates a new theme manager with terminal background detection
func New() *ThemeManager {
	bgType := DetectTerminalBackground()
	tm := &ThemeManager{
		backgroundType: bgType,
		chromaTheme:    getDefaultChromaTheme(bgType),
	}
	tm.colors = tm.buildColorScheme(bgType)
	return tm
}

// NewWithBackground creates a new theme manager with explicit background type
func NewWithBackground(bgType BackgroundType) *ThemeManager {
	tm := &ThemeManager{
		backgroundType: bgType,
		chromaTheme:    getDefaultChromaTheme(bgType),
	}
	tm.colors = tm.buildColorScheme(bgType)
	return tm
}

func getDefaultChromaTheme(bgType BackgroundType) string {
	switch bgType {
	case BackgroundLight:
		return "github" // Light theme for light backgrounds
	case BackgroundDark:
		return "monokai" // Dark theme for dark backgrounds (v1 default)
	default:
		return "monokai" // Default to monokai for unknown backgrounds
	}
}

func (tm *ThemeManager) buildColorScheme(bgType BackgroundType) map[string]string {
	switch bgType {
	case BackgroundLight:
		return tm.buildLightColorScheme()
	case BackgroundDark:
		return tm.buildDarkColorScheme()
	default:
		return tm.buildDarkColorScheme() // Default to dark
	}
}

func (tm *ThemeManager) buildDarkColorScheme() map[string]string {
	return map[string]string{
		string(Header1):       "\033[1;96m",     // Bold Bright Cyan
		string(Header2):       "\033[1;94m",     // Bold Bright Blue
		string(Header3):       "\033[1;95m",     // Bold Bright Magenta
		string(Header4):       "\033[1;93m",     // Bold Bright Yellow
		string(Header5):       "\033[1;92m",     // Bold Bright Green
		string(Header6):       "\033[1;91m",     // Bold Bright Red
		string(Bold):          "\033[1m",        // Bold
		string(Italic):        "\033[3m",        // Italic
		string(Strikethrough): "\033[9m",        // Strikethrough
		string(Code):          "\033[38;5;208m", // Orange (256-color)
		string(BlockQuote):    "\033[38;5;244m", // Gray (256-color)
		string(Link):          "\033[4;94m",     // Underlined Bright Blue
		string(BulletPoint):   "\033[1;97m",     // Bold Bright White
		string(OrderedList):   "\033[1;97m",     // Bold Bright White
		string(TableHeader):   "\033[1;97m",     // Bold Bright White
		string(TableBorder):   "\033[38;5;244m", // Gray (256-color)
		string(Reset):         "\033[0m",        // Reset
	}
}

func (tm *ThemeManager) buildLightColorScheme() map[string]string {
	return map[string]string{
		string(Header1):       "\033[1;34m",     // Bold Blue
		string(Header2):       "\033[1;36m",     // Bold Cyan
		string(Header3):       "\033[1;35m",     // Bold Magenta
		string(Header4):       "\033[1;33m",     // Bold Yellow
		string(Header5):       "\033[1;32m",     // Bold Green
		string(Header6):       "\033[1;31m",     // Bold Red
		string(Bold):          "\033[1m",        // Bold
		string(Italic):        "\033[3m",        // Italic
		string(Strikethrough): "\033[9m",        // Strikethrough
		string(Code):          "\033[38;5;166m", // Dark Orange (256-color)
		string(BlockQuote):    "\033[38;5;240m", // Dark Gray (256-color)
		string(Link):          "\033[4;34m",     // Underlined Blue
		string(BulletPoint):   "\033[1;30m",     // Bold Black
		string(OrderedList):   "\033[1;30m",     // Bold Black
		string(TableHeader):   "\033[1;30m",     // Bold Black
		string(TableBorder):   "\033[38;5;240m", // Dark Gray (256-color)
		string(Reset):         "\033[0m",        // Reset
	}
}

func (tm *ThemeManager) GetColor(key ColorKey) string {
	if color, exists := tm.colors[string(key)]; exists {
		return color
	}
	return tm.colors[string(Reset)]
}

func (tm *ThemeManager) Style(text string, key ColorKey) string {
	return fmt.Sprintf("%s%s%s", tm.GetColor(key), text, tm.GetColor(Reset))
}

func (tm *ThemeManager) StyleNoReset(text string, key ColorKey) string {
	return fmt.Sprintf("%s%s", tm.GetColor(key), text)
}

func (tm *ThemeManager) Reset() string {
	return tm.GetColor(Reset)
}

// GetChromaTheme returns the chroma theme name for code highlighting
// This is used by the CodeHighlighter to apply consistent syntax highlighting
func (tm *ThemeManager) GetChromaTheme() string {
	return tm.chromaTheme
}

// GetTerminalTheme returns detailed terminal theme information
func (tm *ThemeManager) GetTerminalTheme() TerminalTheme {
	return TerminalTheme{
		Background:     tm.backgroundType,
		ChromaTheme:    tm.chromaTheme,
		ColorScheme:    tm.getColorSchemeName(),
		IsHighContrast: tm.isHighContrast(),
	}
}

// SetChromaTheme allows overriding the default chroma theme
func (tm *ThemeManager) SetChromaTheme(themeName string) {
	tm.chromaTheme = themeName
}

// GetBackgroundType returns the detected background type
func (tm *ThemeManager) GetBackgroundType() BackgroundType {
	return tm.backgroundType
}

// IsLightBackground returns true if the background is light
func (tm *ThemeManager) IsLightBackground() bool {
	return tm.backgroundType == BackgroundLight
}

// IsDarkBackground returns true if the background is dark
func (tm *ThemeManager) IsDarkBackground() bool {
	return tm.backgroundType == BackgroundDark
}

func (tm *ThemeManager) getColorSchemeName() string {
	switch tm.backgroundType {
	case BackgroundLight:
		return "light"
	case BackgroundDark:
		return "dark"
	default:
		return "auto"
	}
}

func (tm *ThemeManager) isHighContrast() bool {
	// For v1, we don't have high contrast variants
	// This can be extended in future versions
	return false
}

func (tm *ThemeManager) GetANSIConstant(name string) string {
	constants := map[string]string{
		"reset":          ANSIReset,
		"bold":           ANSIBold,
		"dim":            ANSIDim,
		"italic":         ANSIItalic,
		"underline":      ANSIUnderline,
		"strikethrough":  ANSIStrikethrough,
		"black":          ANSIBlack,
		"red":            ANSIRed,
		"green":          ANSIGreen,
		"yellow":         ANSIYellow,
		"blue":           ANSIBlue,
		"magenta":        ANSIMagenta,
		"cyan":           ANSICyan,
		"white":          ANSIWhite,
		"bright_black":   ANSIBrightBlack,
		"bright_red":     ANSIBrightRed,
		"bright_green":   ANSIBrightGreen,
		"bright_yellow":  ANSIBrightYellow,
		"bright_blue":    ANSIBrightBlue,
		"bright_magenta": ANSIBrightMagenta,
		"bright_cyan":    ANSIBrightCyan,
		"bright_white":   ANSIBrightWhite,
	}

	if constant, exists := constants[name]; exists {
		return constant
	}
	return ANSIReset
}

func (tm *ThemeManager) SupportsColor() bool {
	// Check common environment variables that indicate color support
	term, exists := os.LookupEnv("TERM")
	if !exists {
		return false
	}

	// Check for color support indicators
	if term == "dumb" {
		return false
	}

	// Check for COLORTERM environment variable
	if colorterm, exists := os.LookupEnv("COLORTERM"); exists {
		if colorterm == "truecolor" || colorterm == "24bit" {
			return true
		}
	}

	// Check for common color-supporting terminals
	colorTerms := []string{
		"xterm", "xterm-256color", "xterm-color",
		"screen", "screen-256color",
		"tmux", "tmux-256color",
		"rxvt", "rxvt-unicode",
		"alacritty", "kitty", "wezterm",
	}

	for _, colorTerm := range colorTerms {
		if strings.Contains(term, colorTerm) {
			return true
		}
	}

	return true // Default to supporting color
}

