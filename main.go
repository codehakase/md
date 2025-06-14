package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/codehakase/md/internal/highlighter"
	"github.com/codehakase/md/internal/renderer"
	"github.com/codehakase/md/internal/theme"
	"github.com/codehakase/md/internal/viewer"
)

var (
	vimMode   bool
	watchMode bool
)

var rootCmd = &cobra.Command{
	Use:   "md [flags] <markdown-file>",
	Short: "A markdown renderer and viewer for the terminal",
	Long: `md is a command-line tool that renders markdown files with syntax highlighting
and provides options for vim-style navigation.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]
		if !filepath.IsAbs(filename) {
			var err error
			filename, err = filepath.Abs(filename)
			if err != nil {
				return fmt.Errorf("error resolving file path: %v", err)
			}
		}

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", filename)
		}

		themeManager := theme.New()
		mdRenderer := renderer.New(themeManager)
		codeHighlighter := highlighter.New(themeManager)
		mdViewer := viewer.New()

		renderAndDisplay := func() error {
			content, err := mdRenderer.RenderFile(filename, codeHighlighter)
			if err != nil {
				return fmt.Errorf("rendering error: %v", err)
			}

			if vimMode {
				return mdViewer.DisplayInVimMode(content)
			} else {
				fmt.Print(content)
				return nil
			}
		}

		return renderAndDisplay()
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&vimMode, "vim", "v", false, "Enable vim-style navigation")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
