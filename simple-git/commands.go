// Core Git command implementations
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Initialize a new Git repository
func initRepo() error {
	// Create .git directory
	if err := os.MkdirAll(".git", 0755); err != nil {
		return err
	}

	// Create subdirectories
	dirs := []string{
		".git/objects",
		".git/refs",
		".git/refs/heads",
		".git/refs/tags",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Create HEAD file pointing to master branch
	headContent := "ref: refs/heads/master\n"
	if err := os.WriteFile(".git/HEAD", []byte(headContent), 0644); err != nil {
		return err
	}

	// Create config file
	configContent := "[core]\n\trepositoryformatversion = 0\n\tfilemode = true\n\tbare = false\n"
	if err := os.WriteFile(".git/config", []byte(configContent), 0644); err != nil {
		return err
	}

	// Create description file
	descContent := "Unnamed repository; edit this file 'description' to name the repository.\n"
	if err := os.WriteFile(".git/description", []byte(descContent), 0644); err != nil {
		return err
	}

	fmt.Println("Initialized empty Git repository in .git/")
	return nil
}

// Add file(s) to the staging area
func addFiles(paths []string) error {
	// Load current index
	idx, err := loadIndex()
	if err != nil {
		return err
	}

	// Process each path
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("cannot add '%s': %v", path, err)
		}

		if info.IsDir() {
			// Add all files in directory
			err := filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// Skip directories and .git
				if fileInfo.IsDir() {
					if fileInfo.Name() == ".git" {
						return filepath.SkipDir
					}
					return nil
				}

				// Add file to index
				return idx.AddFile(filePath)
			})

			if err != nil {
				return err
			}
		} else {
			// Add single file to index
			if err := idx.AddFile(path); err != nil {
				return err
			}
		}
	}

	fmt.Printf("Added %d file(s) to index\n", len(paths))
	return nil
}

// Create a new commit with the current index
func commitChanges(message string) error {
	if message == "" {
		return fmt.Errorf("empty commit message")
	}

	// Load index
	idx, err := loadIndex()
	if err != nil {
		return err
	}

	// Check if there are changes to commit
	if len(idx.Entries) == 0 {
		return fmt.Errorf("nothing to commit (empty index)")
	}

	// Write index as tree
	treeHash, err := idx.WriteTree()
	if err != nil {
		return err
	}

	// Get current HEAD commit as parent
	parentCommit, branch, err := getHEAD()
	if err != nil {
		return err
	}

	// Get author/committer info (normally from config or environment)
	// For simplicity, we'll use hardcoded values
	name := os.Getenv("GIT_AUTHOR_NAME")
	email := os.Getenv("GIT_AUTHOR_EMAIL")
	if name == "" {
		name = "Simple Git User"
	}
	if email == "" {
		email = "user@example.com"
	}

	// Format author/committer strings with timestamp
	timestamp := time.Now()
	timezone := timestamp.Format("-0700")
	authorInfo := fmt.Sprintf("%s <%s> %d %s", name, email, timestamp.Unix(), timezone)
	committerInfo := authorInfo // Same for simplicity

	// Create commit object
	commitHash, err := createCommit(treeHash, parentCommit, message, authorInfo, committerInfo)
	if err != nil {
		return err
	}

	// Update branch pointer
	if branch == "" {
		// Default to master if no branch exists
		branch = DefaultBranch
	}
	if err := updateHEAD(commitHash, branch); err != nil {
		return err
	}

	fmt.Printf("[%s %s] %s\n", branch, commitHash[:7], message)
	return nil
}

// Show the status of the working directory
func showStatus() error {
	// Get current branch
	_, branch, err := getHEAD()
	if err != nil {
		return err
	}

	if branch == "" {
		branch = "No branch"
	}

	fmt.Printf("On branch %s\n", branch)

	// Load index
	idx, err := loadIndex()
	if err != nil {
		return err
	}

	// Track modified, staged, and untracked files
	var stagedFiles []string
	var modifiedFiles []string
	var untrackedFiles []string

	// Check for modified/staged files
	for _, entry := range idx.Entries {
		_, err := os.Stat(entry.Path)
		if os.IsNotExist(err) {
			// File in index but deleted in working directory
			modifiedFiles = append(modifiedFiles, entry.Path+" (deleted)")
			continue
		}

		modified, err := idx.IsModified(entry.Path)
		if err != nil {
			return err
		}

		if modified {
			modifiedFiles = append(modifiedFiles, entry.Path)
		}
	}

	// Find untracked files
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and .git
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file is in index
		var found bool
		for _, entry := range idx.Entries {
			if entry.Path == path {
				found = true
				break
			}
		}

		if !found {
			untrackedFiles = append(untrackedFiles, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Print status
	if len(idx.Entries) == 0 {
		fmt.Println("No commits yet")
	}

	if len(modifiedFiles) == 0 && len(stagedFiles) == 0 {
		fmt.Println("No changes to commit (working directory clean)")
	} else {
		if len(modifiedFiles) > 0 {
			fmt.Println("\nChanges not staged for commit:")
			for _, file := range modifiedFiles {
				fmt.Printf("  modified: %s\n", file)
			}
			fmt.Println("\nUse \"git add <file>...\" to update what will be committed")
		}

		if len(stagedFiles) > 0 {
			fmt.Println("\nChanges to be committed:")
			for _, file := range stagedFiles {
				fmt.Printf("  new file: %s\n", file)
			}
		}
	}

	if len(untrackedFiles) > 0 {
		fmt.Println("\nUntracked files:")
		for _, file := range untrackedFiles {
			fmt.Printf("  %s\n", file)
		}
		fmt.Println("\nUse \"git add <file>...\" to include in what will be committed")
	}

	return nil
}

// Show commit history
func showLog() error {
	// Get current HEAD
	headCommit, _, err := getHEAD()
	if err != nil {
		return err
	}

	if headCommit == "" {
		fmt.Println("No commits yet")
		return nil
	}

	// Walk commit chain
	currentCommit := headCommit
	for currentCommit != "" {
		// Read commit object
		commitObj, err := readObject(currentCommit)
		if err != nil {
			return err
		}

		if commitObj.Type != CommitType {
			return fmt.Errorf("not a commit object: %s", currentCommit)
		}

		// Parse commit data
		lines := strings.Split(string(commitObj.Data), "\n")
		var author, message, parent string

		// Find commit info
		messageStart := false
		for _, line := range lines {
			if messageStart {
				if message == "" {
					message = line
				} else {
					message += "\n    " + line
				}
			} else if line == "" {
				// Empty line separates headers from message
				messageStart = true
			} else if strings.HasPrefix(line, "author ") {
				author = strings.TrimPrefix(line, "author ")
			} else if strings.HasPrefix(line, "tree ") {
				continue
			} else if strings.HasPrefix(line, "parent ") {
				parent = strings.TrimPrefix(line, "parent ")
			}
		}

		// Extract author name and date
		authorParts := strings.Split(author, " ")
		authorName := strings.Join(authorParts[:len(authorParts)-2], " ")
		authorTime := authorParts[len(authorParts)-2]

		// Format output
		fmt.Printf("commit %s\n", currentCommit)
		fmt.Printf("Author: %s\n", authorName)
		fmt.Printf("Date:   %s\n", formatUnixTime(authorTime))
		fmt.Printf("\n    %s\n\n", message)

		// Move to parent commit
		currentCommit = parent
	}

	return nil
}

// Format Unix timestamp as readable date
func formatUnixTime(timestamp string) string {
	// Simple implementation - in a real Git, this would use proper date formatting
	return timestamp
}

// Create a new branch
func createNewBranch(name string) error {
	if name == "" {
		return fmt.Errorf("branch name required")
	}

	if err := createBranch(name); err != nil {
		return err
	}

	fmt.Printf("Created branch '%s'\n", name)
	return nil
}

// List all branches
func listAllBranches() error {
	// Get current branch
	_, currentBranch, err := getHEAD()
	if err != nil {
		return err
	}

	if currentBranch == "" {
		currentBranch = "master" // Default
	}

	// Get all branches
	branches, err := listBranches()
	if err != nil {
		return err
	}

	// Handle case where no branches exist yet
	if len(branches) == 0 {
		fmt.Println("No branches")
		return nil
	}

	// Print branches
	for _, branch := range branches {
		if branch == currentBranch {
			fmt.Printf("* %s\n", branch)
		} else {
			fmt.Printf("  %s\n", branch)
		}
	}

	return nil
}

// Checkout a branch
func checkoutBranchOrCommit(name string) error {
	if name == "" {
		return fmt.Errorf("branch name required")
	}

	// Check if it's a branch name
	branchPath := filepath.Join(".git", "refs", "heads", name)
	if _, err := os.Stat(branchPath); err == nil {
		// It's a branch
		if err := checkoutBranch(name); err != nil {
			return err
		}

		fmt.Printf("Switched to branch '%s'\n", name)
		return nil
	}

	// If not a branch, check if it's a valid commit
	_, err := readObject(name)
	if err == nil {
		// It's a commit - detached HEAD state
		// This is simplified and doesn't actually check file types
		return fmt.Errorf("detached HEAD checkout not implemented")
	}

	return fmt.Errorf("not a valid branch or commit: %s", name)
}
