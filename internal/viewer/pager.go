package viewer

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

// Pager handles the integration with the less command for displaying content
type Pager struct {
	lessPath string
}

// NewPager creates a new Pager instance
func NewPager() *Pager {
	return &Pager{
		lessPath: findLessCommand(),
	}
}

// IsLessAvailable checks if the less command is available on the system
func (p *Pager) IsLessAvailable() bool {
	return p.lessPath != ""
}

// Display shows the content using less with vim-style navigation options
func (p *Pager) Display(content string) error {
	if p.lessPath == "" {
		return fmt.Errorf("less command not available")
	}

	args := []string{
		"-R", // Raw control characters (for ANSI colors)
		"-S", // Chop long lines (don't wrap)
		"-X", // Don't clear screen on exit
		"-F", // Quit if entire file fits on screen
		"-K", // Exit on Ctrl-C
		"+g", // Start at beginning (gg equivalent)
	}

	env := append(os.Environ(),
		"LESS_TERMCAP_md=\033[1;36m",    // Bold cyan for headings
		"LESS_TERMCAP_us=\033[1;32m",    // Bold green for underline
		"LESS_TERMCAP_so=\033[1;44;33m", // Bold yellow on blue for standout
		"LESS_TERMCAP_se=\033[0m",       // End standout
		"LESS_TERMCAP_ue=\033[0m",       // End underline
		"LESS_TERMCAP_me=\033[0m",       // End bold/italic
	)

	cmd := exec.Command(p.lessPath, args...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		stdin.Close()
		return fmt.Errorf("failed to start less: %w", err)
	}

	go func() {
		defer stdin.Close()
		stdin.Write([]byte(content))
	}()

	err = cmd.Wait()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() <= 1 {
					return nil
				}
			}
		}
		return fmt.Errorf("less command failed: %w", err)
	}

	return nil
}

// Close performs cleanup (currently no resources to clean up)
func (p *Pager) Close() error {
	return nil
}

// findLessCommand attempts to locate the less command on the system
// TODO (codehakase): expand runtime checks, current version is non deterministic
func findLessCommand() string {
	candidates := []string{"less"}

	if runtime.GOOS == "darwin" {
		candidates = append(candidates,
			"/opt/homebrew/bin/less", // Apple Silicon Homebrew
			"/usr/local/bin/less",    // Intel Homebrew
		)
	}

	if runtime.GOOS == "linux" {
		candidates = append(candidates,
			"/usr/bin/less",
			"/bin/less",
		)
	}

	for _, candidate := range candidates {
		if path, err := exec.LookPath(candidate); err == nil {
			return path
		}
	}

	return ""
}
