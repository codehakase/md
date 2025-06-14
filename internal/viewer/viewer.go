package viewer

import (
	"fmt"
	"runtime"
	"strings"
)

// Viewer provides a vim-style interface for viewing rendered markdown content
type Viewer struct {
	pager *Pager
}

// New creates a new Viewer instance
func New() *Viewer {
	return &Viewer{
		pager: NewPager(),
	}
}

// DisplayInVimMode displays the given content in a less-like interface with vim-style navigation
func (v *Viewer) DisplayInVimMode(content string) error {
	if content == "" {
		return fmt.Errorf("no content to display")
	}

	if !v.pager.IsLessAvailable() {
		return v.fallbackDisplay(content)
	}

	return v.pager.Display(content)
}

// fallbackDisplay provides a simple fallback when less is not available
func (v *Viewer) fallbackDisplay(content string) error {
	fmt.Println("Note: 'less' command not available, displaying content directly:")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Print(content)
	if !strings.HasSuffix(content, "\n") {
		fmt.Println()
	}
	fmt.Println(strings.Repeat("-", 80))
	
	if runtime.GOOS == "windows" {
		fmt.Print("Press Enter to continue...")
		fmt.Scanln()
	}
	
	return nil
}

// Close performs cleanup when the viewer is no longer needed
func (v *Viewer) Close() error {
	if v.pager != nil {
		return v.pager.Close()
	}
	return nil
}