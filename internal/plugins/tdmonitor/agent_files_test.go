package tdmonitor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveAgentFile(t *testing.T) {
	tests := []struct {
		name     string
		files    []string // files to create in temp dir
		expected string   // expected filename (not full path)
	}{
		{
			name:     "prefers AGENTS.md when it exists",
			files:    []string{"AGENTS.md"},
			expected: "AGENTS.md",
		},
		{
			name:     "falls back to CLAUDE.md when no AGENTS.md",
			files:    []string{"CLAUDE.md"},
			expected: "CLAUDE.md",
		},
		{
			name:     "falls back to GEMINI.md when no AGENTS.md or CLAUDE.md",
			files:    []string{"GEMINI.md"},
			expected: "GEMINI.md",
		},
		{
			name:     "AGENTS.md wins over CLAUDE.md",
			files:    []string{"AGENTS.md", "CLAUDE.md"},
			expected: "AGENTS.md",
		},
		{
			name:     "AGENTS.md wins over GEMINI.md",
			files:    []string{"AGENTS.md", "GEMINI.md"},
			expected: "AGENTS.md",
		},
		{
			name:     "AGENTS.md wins over all",
			files:    []string{"AGENTS.md", "CLAUDE.md", "GEMINI.md"},
			expected: "AGENTS.md",
		},
		{
			name:     "CLAUDE.md wins over GEMINI.md",
			files:    []string{"CLAUDE.md", "GEMINI.md"},
			expected: "CLAUDE.md",
		},
		{
			name:     "defaults to AGENTS.md when nothing exists",
			files:    nil,
			expected: "AGENTS.md",
		},
		{
			name:     "ignores local files and defaults to AGENTS.md",
			files:    []string{"CLAUDE.local.md", "GEMINI.local.md"},
			expected: "AGENTS.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			for _, f := range tt.files {
				if err := os.WriteFile(filepath.Join(tmpDir, f), []byte("# Instructions"), 0644); err != nil {
					t.Fatalf("failed to create %s: %v", f, err)
				}
			}

			got := resolveAgentFile(tmpDir)
			expected := filepath.Join(tmpDir, tt.expected)
			if got != expected {
				t.Errorf("resolveAgentFile() = %q, want %q", got, expected)
			}
		})
	}
}
