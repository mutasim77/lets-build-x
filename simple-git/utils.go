// Helper functions and utilities
package main

import (
	"fmt"
	"os"
)

// Check if a directory is a Git repository
func isGitRepo() bool {
	gitDir := ".git"
	info, err := os.Stat(gitDir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Error handling helper - prints error and exits
func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}
