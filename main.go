package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/codehakase/md/internal/highlighter"
	"github.com/codehakase/md/internal/renderer"
	"github.com/codehakase/md/internal/theme"
	"github.com/codehakase/md/internal/viewer"
)

var (
	plainMode bool
)

var rootCmd = &cobra.Command{
	Use:   "md [flags] [markdown-file]",
	Short: "A markdown renderer and viewer for the terminal",
	Long: `md is a command-line tool that renders markdown files with syntax highlighting
and provides options for vim-style navigation.

If no file is provided, md reads from stdin.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		themeManager := theme.New()
		mdRenderer := renderer.New(themeManager)
		codeHighlighter := highlighter.New(themeManager)
		mdViewer := viewer.New()

		var content string
		var err error

		if len(args) == 0 || args[0] == "-" {
			// Read from stdin
			stdinContent, readErr := io.ReadAll(os.Stdin)
			if readErr != nil {
				return fmt.Errorf("error reading from stdin: %v", readErr)
			}
			if len(stdinContent) == 0 {
				return fmt.Errorf("no input provided")
			}
			content, err = mdRenderer.RenderContent(stdinContent, codeHighlighter)
		} else {
			filename := args[0]
			if !filepath.IsAbs(filename) {
				filename, err = filepath.Abs(filename)
				if err != nil {
					return fmt.Errorf("error resolving file path: %v", err)
				}
			}

			if _, statErr := os.Stat(filename); os.IsNotExist(statErr) {
				return fmt.Errorf("file not found: %s", filename)
			}

			content, err = mdRenderer.RenderFile(filename, codeHighlighter)
		}

		if err != nil {
			return fmt.Errorf("rendering error: %v", err)
		}

		if plainMode {
			fmt.Print(content)
			return nil
		}
		return mdViewer.DisplayInVimMode(content)
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&plainMode, "plain", "p", false, "Render entire markdown")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
