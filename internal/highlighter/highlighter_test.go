package highlighter

import (
	"strings"
	"testing"

	"github.com/codehakase/md/internal/theme"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tm := theme.New()
	h := New(tm)

	if h == nil {
		t.Fatal("New() returned nil")
	}

	if h.themeManager != tm {
		t.Error("themeManager not set correctly")
	}

	if h.chromaHelper == nil {
		t.Error("chromaHelper not initialized")
	}
}

func TestHighlightCode(t *testing.T) {
	t.Parallel()

	tm := theme.New()
	h := New(tm)

	tests := []struct {
		name     string
		code     string
		language string
		wantErr  bool
	}{
		{
			name:     "Go code",
			code:     `package main\n\nfunc main() {\n\tfmt.Println("Hello, World!")\n}`,
			language: "go",
			wantErr:  false,
		},
		{
			name:     "Python code",
			code:     `def hello():\n    print("Hello, World!")`,
			language: "python",
			wantErr:  false,
		},
		{
			name:     "JavaScript code",
			code:     `function hello() {\n    console.log("Hello, World!");\n}`,
			language: "javascript",
			wantErr:  false,
		},
		{
			name:     "Empty code",
			code:     "",
			language: "go",
			wantErr:  false,
		},
		{
			name:     "Whitespace only",
			code:     "   \n\t  \n   ",
			language: "go",
			wantErr:  false,
		},
		{
			name:     "Unknown language",
			code:     `some code here`,
			language: "unknownlang",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := h.HighlightCode(tt.code, tt.language)

			if (err != nil) != tt.wantErr {
				t.Errorf("HighlightCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if strings.TrimSpace(tt.code) == "" {
				// For empty/whitespace-only code, result should be the same
				if result != tt.code {
					t.Errorf("HighlightCode() for empty code = %v, want %v", result, tt.code)
				}
			} else {
				// For non-empty code, result should not be empty
				if result == "" {
					t.Error("HighlightCode() returned empty result for non-empty code")
				}
			}
		})
	}
}

func TestHighlightInlineCode(t *testing.T) {
	t.Parallel()

	tm := theme.New()
	h := New(tm)

	code := "fmt.Println()"
	result := h.HighlightInlineCode(code)

	// Should contain ANSI codes from theme manager
	if !strings.Contains(result, "\033[") {
		t.Error("HighlightInlineCode() should contain ANSI escape codes")
	}

	// Should contain the original code
	if !strings.Contains(result, code) {
		t.Error("HighlightInlineCode() should contain the original code")
	}
}

func TestNormalizeLanguage(t *testing.T) {
	t.Parallel()

	tm := theme.New()
	h := New(tm)

	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"go", "go"},
		{"golang", "go"},
		{"js", "javascript"},
		{"javascript", "javascript"},
		{"py", "python"},
		{"python", "python"},
		{"ts", "typescript"},
		{"typescript", "typescript"},
		{"sh", "bash"},
		{"shell", "bash"},
		{"bash", "bash"},
		{"yml", "yaml"},
		{"yaml", "yaml"},
		{"cpp", "cpp"},
		{"c++", "cpp"},
		{"cxx", "cpp"},
		{"dockerfile", "dockerfile"},
		{"docker", "dockerfile"},
		{"plain", "text"},
		{"text", "text"},
		{"txt", "text"},
		{"PYTHON", "python"},
		{"  go  ", "go"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			result := h.normalizeLanguage(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeLanguage(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

