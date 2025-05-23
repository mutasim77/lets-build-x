// Git references (branches, tags, HEAD) implementation
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Default branch name
const DefaultBranch = "master"

// Get the commit that HEAD is pointing to
// Returns the commit hash and branch name (if on a branch)
func getHEAD() (string, string, error) {
	// Read HEAD file
	headPath := filepath.Join(".git", "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return "", "", err
	}

	headString := strings.TrimSpace(string(headContent))

	// Check if HEAD is a reference or direct commit
	if strings.HasPrefix(headString, "ref: ") {
		// HEAD points to a reference
		refPath := strings.TrimPrefix(headString, "ref: ")

		// Extract branch name from reference
		var branchName string
		if strings.HasPrefix(refPath, "refs/heads/") {
			branchName = strings.TrimPrefix(refPath, "refs/heads/")
		}

		// Read the reference to get commit hash
		commitHash, err := readRef(refPath)
		if err != nil {
			// Reference might not exist yet
			if os.IsNotExist(err) {
				return "", branchName, nil
			}
			return "", "", err
		}

		return commitHash, branchName, nil
	} else {
		// HEAD points directly to a commit (detached HEAD)
		return headString, "", nil
	}
}

// Update HEAD to point to a specific commit
// If branch is specified, also update that branch
func updateHEAD(commitHash string, branch string) error {
	if branch != "" {
		// Update HEAD to point to branch
		headContent := fmt.Sprintf("ref: refs/heads/%s", branch)
		headPath := filepath.Join(".git", "HEAD")
		if err := os.WriteFile(headPath, []byte(headContent), 0644); err != nil {
			return err
		}

		// Update branch to point to commit
		return updateRef(fmt.Sprintf("refs/heads/%s", branch), commitHash)
	} else {
		// Detached HEAD - point directly to commit
		headPath := filepath.Join(".git", "HEAD")
		return os.WriteFile(headPath, []byte(commitHash), 0644)
	}
}

// Read a reference (branch, tag) and return the commit hash it points to
func readRef(refPath string) (string, error) {
	fullPath := filepath.Join(".git", refPath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

// Update a reference to point to a specific commit
func updateRef(refPath string, commitHash string) error {
	fullPath := filepath.Join(".git", refPath)

	// Create directory if needed
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write commit hash to reference file
	return os.WriteFile(fullPath, []byte(commitHash), 0644)
}

// Create a new branch pointing to the current HEAD
func createBranch(branchName string) error {
	// Check if branch already exists
	branchPath := filepath.Join(".git", "refs", "heads", branchName)
	if _, err := os.Stat(branchPath); err == nil {
		return fmt.Errorf("branch already exists: %s", branchName)
	}

	// Get current HEAD commit
	headCommit, _, err := getHEAD()
	if err != nil {
		return err
	}

	if headCommit == "" {
		return fmt.Errorf("cannot create branch: no commits yet")
	}

	// Create branch pointing to current HEAD
	return updateRef(fmt.Sprintf("refs/heads/%s", branchName), headCommit)
}

// List all branches
func listBranches() ([]string, error) {
	branches := []string{}

	// Read all files in refs/heads directory
	headsDir := filepath.Join(".git", "refs", "heads")
	err := filepath.Walk(headsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Extract branch name from path
		relPath, err := filepath.Rel(headsDir, path)
		if err != nil {
			return err
		}

		branches = append(branches, relPath)
		return nil
	})

	// Handle case where refs/heads doesn't exist yet
	if os.IsNotExist(err) {
		return branches, nil
	}

	return branches, err
}

// Change HEAD to point to a different branch
func checkoutBranch(branchName string) error {
	// Check if branch exists
	branchPath := filepath.Join(".git", "refs", "heads", branchName)
	if _, err := os.Stat(branchPath); os.IsNotExist(err) {
		return fmt.Errorf("branch does not exist: %s", branchName)
	}

	// Get the commit hash the branch points to
	commitHash, err := readRef(fmt.Sprintf("refs/heads/%s", branchName))
	if err != nil {
		return err
	}

	// Update HEAD to point to branch
	if err := updateHEAD(commitHash, branchName); err != nil {
		return err
	}

	// Read commit to get tree
	commit, err := readObject(commitHash)
	if err != nil {
		return err
	}

	// Parse commit object to get tree hash
	lines := strings.Split(string(commit.Data), "\n")
	var treeHash string
	for _, line := range lines {
		if strings.HasPrefix(line, "tree ") {
			treeHash = strings.TrimPrefix(line, "tree ")
			break
		}
	}

	if treeHash == "" {
		return fmt.Errorf("invalid commit object: no tree found")
	}

	// Rebuild working directory from tree
	return checkoutTree(treeHash, "")
}

// Recursively checkout a tree object to the working directory
func checkoutTree(treeHash string, prefix string) error {
	// Read tree object
	treeObj, err := readObject(treeHash)
	if err != nil {
		return err
	}

	if treeObj.Type != TreeType {
		return fmt.Errorf("not a tree object: %s", treeHash)
	}

	// Parse tree entries
	entries, err := parseTree(treeObj.Data)
	if err != nil {
		return err
	}

	// Process each entry
	for _, entry := range entries {
		// Use map access instead of field access
		path := filepath.Join(prefix, entry["name"])

		// Check mode to determine type
		mode := entry["mode"]

		if mode == DirectoryMode || strings.HasPrefix(mode, "040") {
			// Subdirectory - create and recurse
			fullPath := filepath.Join(".", path)
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				return err
			}

			if err := checkoutTree(entry["hash"], path); err != nil {
				return err
			}
		} else {
			// Regular file - extract from object database
			blob, err := readObject(entry["hash"])
			if err != nil {
				return err
			}

			if blob.Type != BlobType {
				return fmt.Errorf("not a blob object: %s", entry["hash"])
			}

			// Create parent directories if needed
			fullPath := filepath.Join(".", path)
			parentDir := filepath.Dir(fullPath)
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				return err
			}

			// Write file
			fileMode := 0644
			if mode == ExecutableMode || strings.HasPrefix(mode, "100755") {
				fileMode = 0755
			}

			if err := os.WriteFile(fullPath, blob.Data, os.FileMode(fileMode)); err != nil {
				return err
			}
		}
	}

	return nil
}
