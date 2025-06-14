package highlighter

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

// ChromaHelper handles the low-level Chroma integration for syntax highlighting
type ChromaHelper struct {
	formatter chroma.Formatter
	style     *chroma.Style
}

// NewChromaHelper creates a new ChromaHelper with optimal terminal settings
func NewChromaHelper() *ChromaHelper {
	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Get("terminal")
		if formatter == nil {
			formatter = formatters.Fallback
		}
	}

	style := styles.Get("monokai")
	if style == nil {
		style = styles.Get("dracula")
		if style == nil {
			style = styles.Get("github-dark")
			if style == nil {
				style = styles.Fallback
			}
		}
	}

	return &ChromaHelper{
		formatter: formatter,
		style:     style,
	}
}

// Highlight performs syntax highlighting using Chroma
func (ch *ChromaHelper) Highlight(code, language string) (string, error) {
	lexer := ch.getLexer(language, code)
	if lexer == nil {
		return "", fmt.Errorf("no suitable lexer found for language: %s", language)
	}

	lexer = chroma.Coalesce(lexer)

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return "", fmt.Errorf("failed to tokenize code: %w", err)
	}

	var buf bytes.Buffer
	err = ch.formatter.Format(&buf, ch.style, iterator)
	if err != nil {
		return "", fmt.Errorf("failed to format highlighted code: %w", err)
	}

	result := buf.String()
	result = strings.TrimRight(result, "\n")
	
	return result, nil
}

// getLexer returns the most appropriate lexer for the given language and code
func (ch *ChromaHelper) getLexer(language, code string) chroma.Lexer {
	var lexer chroma.Lexer

	if language != "" {
		lexer = lexers.Get(language)
		if lexer != nil {
			return lexer
		}

		if aliasLexer := ch.getLexerByAlias(language); aliasLexer != nil {
			return aliasLexer
		}
	}

	lexer = lexers.Analyse(code)
	if lexer != nil {
		return lexer
	}

	return lexers.Fallback
}

// getLexerByAlias handles language aliases that might not be directly supported
func (ch *ChromaHelper) getLexerByAlias(language string) chroma.Lexer {
	switch strings.ToLower(language) {
	case "jsonc", "json5":
		return lexers.Get("json")
	case "tsx":
		return lexers.Get("typescript")
	case "jsx":
		return lexers.Get("javascript")
	case "sh", "shell", "zsh", "fish":
		return lexers.Get("bash")
	case "yml":
		return lexers.Get("yaml")
	case "ps1", "powershell":
		return lexers.Get("powershell")
	case "cmd", "batch", "bat":
		return lexers.Get("batch")
	case "asm", "assembly":
		return lexers.Get("nasm")
	case "tex", "latex":
		return lexers.Get("latex")
	case "md", "markdown":
		return lexers.Get("markdown")
	case "ini", "cfg", "conf", "config":
		return lexers.Get("ini")
	case "toml":
		return lexers.Get("toml")
	case "proto", "protobuf":
		return lexers.Get("protobuf")
	case "graphql", "gql":
		return lexers.Get("graphql")
	case "hcl", "terraform":
		return lexers.Get("hcl")
	case "nginx":
		return lexers.Get("nginx")
	case "apache":
		return lexers.Get("apacheconf")
	default:
		return nil
	}
}

// GetAvailableStyles returns a list of available Chroma styles suitable for terminals
func (ch *ChromaHelper) GetAvailableStyles() []string {
	allStyles := styles.Names()
	terminalFriendly := []string{}

	terminalFriendlyNames := map[string]bool{
		"monokai":     true,
		"dracula":     true,
		"github-dark": true,
		"native":      true,
		"fruity":      true,
		"material":    true,
		"nord":        true,
		"onedark":     true,
		"solarized-dark": true,
		"tomorrow-night": true,
		"vs-dark":     true,
	}

	for _, style := range allStyles {
		if terminalFriendlyNames[style] {
			terminalFriendly = append(terminalFriendly, style)
		}
	}

	return terminalFriendly
}

// SetStyle changes the current highlighting style
// Note: chroma always returns a style (fallback if not found), so this never errors
func (ch *ChromaHelper) SetStyle(styleName string) error {
	style := styles.Get(styleName)
	ch.style = style
	return nil
}

// GetCurrentStyle returns the name of the currently active style
func (ch *ChromaHelper) GetCurrentStyle() string {
	if ch.style == nil {
		return "fallback"
	}
	return ch.style.Name
}

