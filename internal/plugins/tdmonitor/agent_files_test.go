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

func TestAnyFileHasTDInstructions(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string // filename -> content
		expected bool
	}{
		{
			name:     "no files exist",
			files:    nil,
			expected: false,
		},
		{
			name:     "AGENTS.md has instructions",
			files:    map[string]string{"AGENTS.md": "run td usage --new-session"},
			expected: true,
		},
		{
			name:     "CLAUDE.md has instructions",
			files:    map[string]string{"CLAUDE.md": "td usage -q"},
			expected: true,
		},
		{
			name:     "CLAUDE.local.md has instructions",
			files:    map[string]string{"CLAUDE.local.md": "td usage --new-session"},
			expected: true,
		},
		{
			name:     "GEMINI.md has instructions",
			files:    map[string]string{"GEMINI.md": "run td usage"},
			expected: true,
		},
		{
			name:     "GEMINI.local.md has instructions",
			files:    map[string]string{"GEMINI.local.md": "td usage -q"},
			expected: true,
		},
		{
			name:     "files exist but no instructions",
			files:    map[string]string{"CLAUDE.md": "# Project\nSome content", "AGENTS.md": "# Agents"},
			expected: false,
		},
		{
			name:     "mixed: one file has instructions among many",
			files:    map[string]string{"AGENTS.md": "# Agents", "CLAUDE.local.md": "td usage --new-session", "GEMINI.md": "# Gemini"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			for name, content := range tt.files {
				if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
					t.Fatalf("failed to create %s: %v", name, err)
				}
			}

			got := anyFileHasTDInstructions(tmpDir)
			if got != tt.expected {
				t.Errorf("anyFileHasTDInstructions() = %v, want %v", got, tt.expected)
			}
		})
	}
}
