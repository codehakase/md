package highlighter

import (
	"strings"

	"github.com/codehakase/md/internal/theme"
)

// Highlighter handles syntax highlighting of code blocks
type Highlighter struct {
	themeManager *theme.ThemeManager
	chromaHelper *ChromaHelper
}

// New creates a new code highlighter with the given theme manager
func New(themeManager *theme.ThemeManager) *Highlighter {
	return &Highlighter{
		themeManager: themeManager,
		chromaHelper: NewChromaHelper(),
	}
}

// Highlight highlights code with the specified language
// It auto-detects language from hints and handles edge cases gracefully
func (h *Highlighter) Highlight(code, language string) (string, error) {
	if strings.TrimSpace(code) == "" {
		return code, nil
	}

	language = h.normalizeLanguage(language)

	highlighted, err := h.chromaHelper.Highlight(code, language)
	if err != nil {
		return h.themeManager.Style(code, theme.Code), nil
	}

	return highlighted, nil
}

// HighlightCode is an alias for Highlight for the interface requirement
func (h *Highlighter) HighlightCode(code, language string) (string, error) {
	return h.Highlight(code, language)
}

// HighlightInlineCode highlights inline code snippets using theme manager
func (h *Highlighter) HighlightInlineCode(code string) string {
	return h.themeManager.Style(code, theme.Code)
}

// normalizeLanguage normalizes language hints from markdown fenced code blocks
func (h *Highlighter) normalizeLanguage(language string) string {
	if language == "" {
		return language
	}

	lang := strings.ToLower(strings.TrimSpace(language))

	switch lang {
	case "js", "javascript":
		return "javascript"
	case "ts", "typescript":
		return "typescript"
	case "py", "python":
		return "python"
	case "rb", "ruby":
		return "ruby"
	case "sh", "bash", "shell":
		return "bash"
	case "yml", "yaml":
		return "yaml"
	case "json":
		return "json"
	case "xml", "html":
		return lang
	case "css":
		return "css"
	case "sql":
		return "sql"
	case "go", "golang":
		return "go"
	case "rust", "rs":
		return "rust"
	case "c":
		return "c"
	case "cpp", "c++", "cxx":
		return "cpp"
	case "java":
		return "java"
	case "php":
		return "php"
	case "swift":
		return "swift"
	case "kotlin", "kt":
		return "kotlin"
	case "scala":
		return "scala"
	case "r":
		return "r"
	case "matlab":
		return "matlab"
	case "perl":
		return "perl"
	case "lua":
		return "lua"
	case "vim":
		return "vim"
	case "dockerfile", "docker":
		return "dockerfile"
	case "makefile", "make":
		return "makefile"
	case "diff":
		return "diff"
	case "plain", "text", "txt":
		return "text"
	default:
		return lang
	}
}