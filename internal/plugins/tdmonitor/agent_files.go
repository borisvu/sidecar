package tdmonitor

import (
	"os"
	"path/filepath"
	"strings"
)

// agentWriteTargets are the non-local agent files eligible as write targets, in priority order.
var agentWriteTargets = []string{"AGENTS.md", "CLAUDE.md", "GEMINI.md"}

// agentFiles is the complete list of agent instruction files to check.
// All files are checked for dedup in anyFileHasTDInstructions.
var agentFiles = []string{
	"AGENTS.md",
	"CLAUDE.md",
	"CLAUDE.local.md",
	"GEMINI.md",
	"GEMINI.local.md",
}

// resolveAgentFile returns the best agent file to install td instructions into.
// Priority: AGENTS.md > CLAUDE.md > GEMINI.md > create AGENTS.md.
// Only tracked/shared files are targets; local files are never written to.
func resolveAgentFile(baseDir string) string {
	for _, name := range agentWriteTargets {
		path := filepath.Join(baseDir, name)
		if fileExists(path) {
			return path
		}
	}
	return filepath.Join(baseDir, "AGENTS.md")
}

// anyFileHasTDInstructions checks all agent files for existing td instructions.
// Returns true if any file contains "td usage", preventing duplicate installation.
func anyFileHasTDInstructions(baseDir string) bool {
	for _, name := range agentFiles {
		path := filepath.Join(baseDir, name)
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		if strings.Contains(string(content), "td usage") {
			return true
		}
	}
	return false
}
