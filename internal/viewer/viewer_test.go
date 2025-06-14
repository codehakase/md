package viewer

import (
	"testing"
)

func TestNew(t *testing.T) {
	t.Parallel()

	v := New()
	if v == nil {
		t.Fatal("New() returned nil")
	}
	if v.pager == nil {
		t.Fatal("New() did not initialize pager")
	}
}

func TestPagerIsLessAvailable(t *testing.T) {
	t.Parallel()

	p := NewPager()
	if p == nil {
		t.Fatal("NewPager() returned nil")
	}

	available := p.IsLessAvailable()
	t.Logf("Less available: %v", available)
}

func TestDisplayInVimModeEmptyContent(t *testing.T) {
	t.Parallel()

	v := New()
	err := v.DisplayInVimMode("")
	if err == nil {
		t.Fatal("Expected error for empty content")
	}
	if err.Error() != "no content to display" {
		t.Fatalf("Expected 'no content to display', got: %v", err)
	}
}

func TestClose(t *testing.T) {
	t.Parallel()

	v := New()
	err := v.Close()
	if err != nil {
		t.Fatalf("Close() returned error: %v", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsHelper(s, substr))))
}

func containsHelper(s, substr string) bool {
	for i := 1; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

