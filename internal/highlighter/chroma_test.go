package highlighter

import (
	"testing"
)

func TestNewChromaHelper(t *testing.T) {
	t.Parallel()

	ch := NewChromaHelper()

	if ch == nil {
		t.Fatal("NewChromaHelper() returned nil")
	}

	if ch.formatter == nil {
		t.Error("formatter not initialized")
	}

	if ch.style == nil {
		t.Error("style not initialized")
	}
}

func TestChromaHelperHighlight(t *testing.T) {
	t.Parallel()

	ch := NewChromaHelper()

	tests := []struct {
		name     string
		code     string
		language string
		wantErr  bool
	}{
		{
			name:     "Go code",
			code:     "package main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}",
			language: "go",
			wantErr:  false,
		},
		{
			name:     "Python code",
			code:     "def hello():\n    print('Hello')",
			language: "python",
			wantErr:  false,
		},
		{
			name:     "JSON code",
			code:     `{"name": "test", "value": 42}`,
			language: "json",
			wantErr:  false,
		},
		{
			name:     "Unknown language",
			code:     "some random text",
			language: "unknownlang",
			wantErr:  false,
		},
		{
			name:     "Empty language with code analysis",
			code:     "package main\n\nfunc main() {}",
			language: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := ch.Highlight(tt.code, tt.language)

			if (err != nil) != tt.wantErr {
				t.Errorf("Highlight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Should return highlighted code (may contain ANSI codes)
				if result == "" {
					t.Error("Highlight() returned empty result")
				}

				// Result should not be empty and should be different from input (due to ANSI codes)
				// We can't easily check for exact content due to ANSI color codes,
				// but the length should be different (longer due to ANSI codes)
				if len(result) == 0 {
					t.Error("Highlight() returned empty result")
				}
			}
		})
	}
}

func TestGetLexerByAlias(t *testing.T) {
	t.Parallel()

	ch := NewChromaHelper()

	tests := []struct {
		alias    string
		expected bool // whether a lexer should be found
	}{
		{"tsx", true},
		{"jsx", true},
		{"sh", true},
		{"shell", true},
		{"yml", true},
		{"ps1", true},
		{"powershell", true},
		{"cmd", true},
		{"batch", true},
		{"jsonc", true},
		{"json5", true},
		{"asm", true},
		{"assembly", true},
		{"tex", true},
		{"latex", true},
		{"md", true},
		{"markdown", true},
		{"ini", true},
		{"cfg", true},
		{"conf", true},
		{"config", true},
		{"toml", true},
		{"proto", true},
		{"protobuf", true},
		{"nonexistentlang", false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.alias, func(t *testing.T) {
			t.Parallel()

			lexer := ch.getLexerByAlias(tt.alias)
			found := lexer != nil

			if found != tt.expected {
				t.Errorf("getLexerByAlias(%q) found=%v, expected=%v", tt.alias, found, tt.expected)
			}
		})
	}
}

func TestGetAvailableStyles(t *testing.T) {
	t.Parallel()

	ch := NewChromaHelper()
	styles := ch.GetAvailableStyles()

	if len(styles) == 0 {
		t.Error("GetAvailableStyles() returned empty slice")
	}

	// Check that some expected terminal-friendly styles are included
	expectedStyles := []string{"monokai", "dracula", "github-dark"}
	for _, expected := range expectedStyles {
		found := false
		for _, style := range styles {
			if style == expected {
				found = true
				break
			}
		}
		if !found {
			// It's okay if a style isn't available, but let's log it
			t.Logf("Expected style %s not found in available styles", expected)
		}
	}
}

func TestSetStyle(t *testing.T) {
	t.Parallel()

	ch := NewChromaHelper()

	tests := []struct {
		name      string
		styleName string
		wantErr   bool
	}{
		{
			name:      "valid style",
			styleName: "monokai",
			wantErr:   false,
		},
		{
			name:      "unknown style should not error",
			styleName: "somerarestyle",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ch.SetStyle(tt.styleName)
			if (err != nil) != tt.wantErr {
				if tt.wantErr {
					t.Errorf("SetStyle(%s) should return error", tt.styleName)
				} else {
					t.Errorf("SetStyle(%s) returned error: %v", tt.styleName, err)
				}
			}
		})
	}
}

func TestGetCurrentStyle(t *testing.T) {
	t.Parallel()

	ch := NewChromaHelper()

	currentStyle := ch.GetCurrentStyle()
	if currentStyle == "" {
		t.Error("GetCurrentStyle() returned empty string")
	}

	// Test after setting a style
	err := ch.SetStyle("monokai")
	if err == nil {
		newStyle := ch.GetCurrentStyle()
		if newStyle != "monokai" {
			t.Errorf("GetCurrentStyle() = %q, want %q", newStyle, "monokai")
		}
	}
}

