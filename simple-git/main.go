// Entry point and CLI handler
package main

import (
	"flag"
	"fmt"
	"os"
)

// Version information
const SimpleGitVersion = "0.1.0"

// Print usage information
func printUsage() {
	fmt.Println("Simple Git - A simplified Git implementation")
	fmt.Println("\nUsage:")
	fmt.Println("  simple-git <command> [<args>]")
	fmt.Println("\nCommands:")
	fmt.Println("  init                Create an empty Git repository")
	fmt.Println("  add <file>...       Add file(s) to the staging area")
	fmt.Println("  commit -m <msg>     Create a new commit")
	fmt.Println("  status              Show the working tree status")
	fmt.Println("  log                 Show commit logs")
	fmt.Println("  branch [<name>]     Create or list branches")
	fmt.Println("  checkout <branch>   Switch branches")
	fmt.Println("  version             Show version information")
	fmt.Println("  help                Show help information")
}

func main() {
	// No arguments provided
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Get command
	command := os.Args[1]

	// Process commands
	switch command {
	case "init":
		err := initRepo()
		if err != nil {
			exitWithError(err)
		}

	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: No files specified")
			os.Exit(1)
		}

		// Check if in a Git repository
		if !isGitRepo() {
			fmt.Println("Error: Not a Git repository")
			os.Exit(1)
		}

		// Add files
		err := addFiles(os.Args[2:])
		if err != nil {
			exitWithError(err)
		}

	case "commit":
		// Check if in a Git repository
		if !isGitRepo() {
			fmt.Println("Error: Not a Git repository")
			os.Exit(1)
		}

		// Parse flags
		commitFlags := flag.NewFlagSet("commit", flag.ExitOnError)
		message := commitFlags.String("m", "", "Commit message")

		// Parse flags starting from os.Args[2]
		commitFlags.Parse(os.Args[2:])

		// Check for commit message
		if *message == "" {
			fmt.Println("Error: Commit message required (-m flag)")
			os.Exit(1)
		}

		// Create commit
		err := commitChanges(*message)
		if err != nil {
			exitWithError(err)
		}

	case "status":
		// Check if in a Git repository
		if !isGitRepo() {
			fmt.Println("Error: Not a Git repository")
			os.Exit(1)
		}

		err := showStatus()
		if err != nil {
			exitWithError(err)
		}

	case "log":
		// Check if in a Git repository
		if !isGitRepo() {
			fmt.Println("Error: Not a Git repository")
			os.Exit(1)
		}

		err := showLog()
		if err != nil {
			exitWithError(err)
		}

	case "branch":
		// Check if in a Git repository
		if !isGitRepo() {
			fmt.Println("Error: Not a Git repository")
			os.Exit(1)
		}

		// If branch name provided, create new branch
		if len(os.Args) > 2 {
			err := createNewBranch(os.Args[2])
			if err != nil {
				exitWithError(err)
			}
		} else {
			// List branches
			err := listAllBranches()
			if err != nil {
				exitWithError(err)
			}
		}

	case "checkout":
		// Check if in a Git repository
		if !isGitRepo() {
			fmt.Println("Error: Not a Git repository")
			os.Exit(1)
		}

		if len(os.Args) < 3 {
			fmt.Println("Error: Branch name required")
			os.Exit(1)
		}

		err := checkoutBranchOrCommit(os.Args[2])
		if err != nil {
			exitWithError(err)
		}

	case "version":
		fmt.Printf("simple-git version %s\n", SimpleGitVersion)

	case "help":
		printUsage()

	default:
		fmt.Printf("Error: Unknown command '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}
