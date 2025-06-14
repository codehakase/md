package renderer

import (
	"strings"
	"unicode/utf8"
)

// ANSI styling constants
const (
	Reset = "\033[0m"
	
	Bold          = "\033[1m"
	Dim           = "\033[2m"
	Italic        = "\033[3m"
	Underline     = "\033[4m"
	Strikethrough = "\033[9m"
	
	// Header prefixes for visual hierarchy
	H1Prefix = "# "
	H2Prefix = "## "
	H3Prefix = "### "
	H4Prefix = "#### "
	H5Prefix = "##### "
	H6Prefix = "###### "
)

// Indent adds indentation to text
func Indent(text string, level int) string {
	if level <= 0 {
		return text
	}
	
	indent := strings.Repeat("  ", level)
	lines := strings.Split(text, "\n")
	
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			lines[i] = indent + line
		}
	}
	
	return strings.Join(lines, "\n")
}

// WrapText wraps text to specified width
func WrapText(text string, width int) string {
	if width <= 0 {
		return text
	}
	
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}
	
	var lines []string
	var currentLine []string
	currentLength := 0
	
	for _, word := range words {
		wordLength := utf8.RuneCountInString(word)
		
		if currentLength > 0 && currentLength+wordLength+1 > width {
			lines = append(lines, strings.Join(currentLine, " "))
			currentLine = []string{word}
			currentLength = wordLength
		} else {
			currentLine = append(currentLine, word)
			if currentLength > 0 {
				currentLength++
			}
			currentLength += wordLength
		}
	}
	
	if len(currentLine) > 0 {
		lines = append(lines, strings.Join(currentLine, " "))
	}
	
	return strings.Join(lines, "\n")
}

// PadRight pads text to specified width
func PadRight(text string, width int) string {
	textWidth := utf8.RuneCountInString(text)
	if textWidth >= width {
		return text
	}
	return text + strings.Repeat(" ", width-textWidth)
}

// TrimTrailingWhitespace removes trailing whitespace from each line
func TrimTrailingWhitespace(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	return strings.Join(lines, "\n")
}

// EnsureTrailingNewline ensures text ends with exactly one newline
func EnsureTrailingNewline(text string) string {
	text = strings.TrimRight(text, "\n")
	return text + "\n"
}