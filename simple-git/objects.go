// Git object database implementation
package main

import (
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Object types in Git
const (
	BlobType   = "blob"
	TreeType   = "tree"
	CommitType = "commit"
)

// Object represents a Git object (blob, tree, or commit)
type Object struct {
	Type string
	Size int
	Data []byte
}

// hash calculates SHA-1 hash of content with header
// This is how Git generates IDs for all objects
func hash(objType string, data []byte) string {
	header := fmt.Sprintf("%s %d\x00", objType, len(data))
	h := sha1.New()
	h.Write([]byte(header))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// Create a new git object from type and data
// Returns the hash ID of the object
func writeObject(objType string, data []byte) (string, error) {
	// Calculate the hash (object ID)
	id := hash(objType, data)

	// Prepare header (type + size + null byte)
	header := fmt.Sprintf("%s %d\x00", objType, len(data))

	// Create the directory structure
	objDir := filepath.Join(".git", "objects", id[:2])
	if err := os.MkdirAll(objDir, 0755); err != nil {
		return "", err
	}

	// Create the object file
	objPath := filepath.Join(objDir, id[2:])
	file, err := os.Create(objPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Set up zlib compression
	// Git compresses all objects for efficiency
	z := zlib.NewWriter(file)
	defer z.Close()

	// Write header + data
	if _, err := z.Write([]byte(header)); err != nil {
		return "", err
	}
	if _, err := z.Write(data); err != nil {
		return "", err
	}

	return id, nil
}

// Read an object from the object database
func readObject(id string) (*Object, error) {
	// Check if ID is valid SHA-1
	if len(id) != 40 {
		return nil, fmt.Errorf("invalid object id: %s", id)
	}

	// Find the object file
	objPath := filepath.Join(".git", "objects", id[:2], id[2:])
	file, err := os.Open(objPath)
	if err != nil {
		return nil, fmt.Errorf("object not found: %s", id)
	}
	defer file.Close()

	// Create zlib reader to decompress
	z, err := zlib.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer z.Close()

	// Read the decompressed content
	content, err := io.ReadAll(z)
	if err != nil {
		return nil, err
	}

	// Parse header (format: "type size\0data")
	header := strings.SplitN(string(content), "\x00", 2)
	if len(header) != 2 {
		return nil, fmt.Errorf("malformed object: %s", id)
	}

	// Extract type and size
	typeAndSize := strings.SplitN(header[0], " ", 2)
	if len(typeAndSize) != 2 {
		return nil, fmt.Errorf("malformed object header: %s", id)
	}

	objType := typeAndSize[0]

	// Create and return the object
	return &Object{
		Type: objType,
		Size: len(header[1]),
		Data: []byte(header[1]),
	}, nil
}

// Create a blob object from file content
func createBlob(data []byte) (string, error) {
	return writeObject(BlobType, data)
}

// Create a tree object from entries
// Each entry has mode, name, and hash
func createTree(entries []byte) (string, error) {
	return writeObject(TreeType, entries)
}

// Create a commit object with tree, parent, and message
func createCommit(tree, parent, message string, author, committer string) (string, error) {
	// Format: tree <hash>\nparent <hash>\nauthor <info>\ncommitter <info>\n\n<message>
	var commitData strings.Builder

	// Add tree reference
	commitData.WriteString(fmt.Sprintf("tree %s\n", tree))

	// Add parent if provided
	if parent != "" {
		commitData.WriteString(fmt.Sprintf("parent %s\n", parent))
	}

	// Add author and committer
	commitData.WriteString(fmt.Sprintf("author %s\n", author))
	commitData.WriteString(fmt.Sprintf("committer %s\n", committer))

	// Add empty line and message
	commitData.WriteString("\n")
	commitData.WriteString(message)

	return writeObject(CommitType, []byte(commitData.String()))
}

// Encodes a tree entry in Git's format
// Format: mode<space>name\0<SHA-1 in binary>
func encodeTreeEntry(mode, name, hash string) ([]byte, error) {
	// Convert hash from hex to binary
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	// Format the entry
	entry := fmt.Sprintf("%s %s", mode, name)
	result := make([]byte, len(entry)+1+len(hashBytes))

	// Copy mode and name
	copy(result, entry)
	// Add null byte separator
	result[len(entry)] = 0
	// Copy binary hash
	copy(result[len(entry)+1:], hashBytes)

	return result, nil
}

// Parse a tree object's data into entries
func parseTree(data []byte) ([]map[string]string, error) {
	entries := []map[string]string{}

	// Keep processing until we've consumed all data
	i := 0
	for i < len(data) {
		// Find the null byte that separates name from hash
		nullPos := -1
		for j := i; j < len(data); j++ {
			if data[j] == 0 {
				nullPos = j
				break
			}
		}

		if nullPos == -1 || nullPos+20 > len(data) {
			return nil, fmt.Errorf("malformed tree object")
		}

		// Parse mode and name (space-separated)
		modeAndName := strings.SplitN(string(data[i:nullPos]), " ", 2)
		if len(modeAndName) != 2 {
			return nil, fmt.Errorf("malformed tree entry")
		}

		// Extract hash (20 bytes after null)
		hash := hex.EncodeToString(data[nullPos+1 : nullPos+21])

		// Add entry to result
		entries = append(entries, map[string]string{
			"mode": modeAndName[0],
			"name": modeAndName[1],
			"hash": hash,
		})

		// Move to next entry
		i = nullPos + 21
	}

	return entries, nil
}
