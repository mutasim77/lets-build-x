package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"
)

// Index file signature and version
const (
	IndexSignature = "DIRC"
	IndexVersion   = 2
)

// FileMode represents file permissions in Git
const (
	FileMode       = "100644" // Regular file
	ExecutableMode = "100755" // Executable file
	SymlinkMode    = "120000" // Symbolic link
	DirectoryMode  = "40000"  // Directory
)

// IndexEntry represents a single file in the index
type IndexEntry struct {
	Ctime time.Time
	Mtime time.Time
	Dev   uint32
	Ino   uint32
	Mode  uint32
	Uid   uint32
	Gid   uint32
	Size  uint32
	Hash  string
	Flags uint16
	Path  string
}

// Index represents Git's staging area
type Index struct {
	Entries []IndexEntry
}

// Load the index file from disk
func loadIndex() (*Index, error) {
	indexPath := filepath.Join(".git", "index")

	// Check if index exists
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		// Return empty index if file doesn't exist
		return &Index{Entries: []IndexEntry{}}, nil
	}

	// Read the index file
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}

	// Check minimum size and signature
	if len(data) < 12 || string(data[0:4]) != IndexSignature {
		return nil, fmt.Errorf("invalid index file format")
	}

	// Parse version and number of entries
	version := binary.BigEndian.Uint32(data[4:8])
	if version != IndexVersion {
		return nil, fmt.Errorf("unsupported index version: %d", version)
	}

	numEntries := binary.BigEndian.Uint32(data[8:12])

	// Create new index
	idx := &Index{
		Entries: make([]IndexEntry, 0, numEntries),
	}

	// Parse entries
	offset := uint32(12) // Start after header
	for i := uint32(0); i < numEntries; i++ {
		// Ensure we have enough data for fixed fields (62 bytes)
		if offset+62 > uint32(len(data)) {
			return nil, fmt.Errorf("truncated index file")
		}

		// Parse entry fields
		entry := IndexEntry{}

		// Timestamps
		ctimeSec := binary.BigEndian.Uint32(data[offset : offset+4])
		ctimeNano := binary.BigEndian.Uint32(data[offset+4 : offset+8])
		entry.Ctime = time.Unix(int64(ctimeSec), int64(ctimeNano))

		mtimeSec := binary.BigEndian.Uint32(data[offset+8 : offset+12])
		mtimeNano := binary.BigEndian.Uint32(data[offset+12 : offset+16])
		entry.Mtime = time.Unix(int64(mtimeSec), int64(mtimeNano))

		// File metadata
		entry.Dev = binary.BigEndian.Uint32(data[offset+16 : offset+20])
		entry.Ino = binary.BigEndian.Uint32(data[offset+20 : offset+24])
		entry.Mode = binary.BigEndian.Uint32(data[offset+24 : offset+28])
		entry.Uid = binary.BigEndian.Uint32(data[offset+28 : offset+32])
		entry.Gid = binary.BigEndian.Uint32(data[offset+32 : offset+36])
		entry.Size = binary.BigEndian.Uint32(data[offset+36 : offset+40])

		// Object hash (20 bytes)
		entry.Hash = fmt.Sprintf("%x", data[offset+40:offset+60])

		// Flags
		entry.Flags = binary.BigEndian.Uint16(data[offset+60 : offset+62])

		// Extract path name (null-terminated)
		pathStart := offset + 62
		pathEnd := pathStart
		for pathEnd < uint32(len(data)) && data[pathEnd] != 0 {
			pathEnd++
		}
		if pathEnd >= uint32(len(data)) {
			return nil, fmt.Errorf("invalid index entry path")
		}
		entry.Path = string(data[pathStart:pathEnd])

		// Move to next entry, align to 8 bytes
		offset = pathEnd + 1
		if offset%8 != 0 {
			offset += 8 - (offset % 8)
		}

		// Add to entries
		idx.Entries = append(idx.Entries, entry)
	}

	return idx, nil
}

// Save the index to disk
func (idx *Index) Save() error {
	// Sort entries by path
	sort.Slice(idx.Entries, func(i, j int) bool {
		return idx.Entries[i].Path < idx.Entries[j].Path
	})

	// Create buffer for index data
	buffer := bytes.NewBuffer(nil)

	// Write header
	binary.Write(buffer, binary.BigEndian, []byte(IndexSignature))
	binary.Write(buffer, binary.BigEndian, uint32(IndexVersion))
	binary.Write(buffer, binary.BigEndian, uint32(len(idx.Entries)))

	// Write entries
	for _, entry := range idx.Entries {
		// Timestamps
		binary.Write(buffer, binary.BigEndian, uint32(entry.Ctime.Unix()))
		binary.Write(buffer, binary.BigEndian, uint32(entry.Ctime.Nanosecond()))
		binary.Write(buffer, binary.BigEndian, uint32(entry.Mtime.Unix()))
		binary.Write(buffer, binary.BigEndian, uint32(entry.Mtime.Nanosecond()))

		// File metadata
		binary.Write(buffer, binary.BigEndian, entry.Dev)
		binary.Write(buffer, binary.BigEndian, entry.Ino)
		binary.Write(buffer, binary.BigEndian, entry.Mode)
		binary.Write(buffer, binary.BigEndian, entry.Uid)
		binary.Write(buffer, binary.BigEndian, entry.Gid)
		binary.Write(buffer, binary.BigEndian, entry.Size)

		// Object hash
		hashBytes, err := hex.DecodeString(entry.Hash)
		if err != nil {
			return err
		}
		buffer.Write(hashBytes)

		// Flags with path length (lower 12 bits)
		flags := entry.Flags & 0xF000            // Keep high 4 bits
		flags |= uint16(len(entry.Path) & 0xFFF) // Set low 12 bits to path length
		binary.Write(buffer, binary.BigEndian, flags)

		// Path with null terminator
		buffer.Write([]byte(entry.Path))
		buffer.WriteByte(0)

		// Padding to 8-byte boundary
		for buffer.Len()%8 != 0 {
			buffer.WriteByte(0)
		}
	}

	// Write index to file
	indexPath := filepath.Join(".git", "index")
	return os.WriteFile(indexPath, buffer.Bytes(), 0644)
}

// Add a file to the index
func (idx *Index) AddFile(path string) error {
	// Get file stats
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Read file content
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Create blob from file content
	hash, err := createBlob(data)
	if err != nil {
		return err
	}

	// Check if file already in index
	for i, entry := range idx.Entries {
		if entry.Path == path {
			// Update existing entry
			idx.Entries[i].Hash = hash
			idx.Entries[i].Size = uint32(len(data))
			idx.Entries[i].Mtime = info.ModTime()
			idx.Entries[i].Ctime = time.Now()
			return idx.Save()
		}
	}

	// Create new entry
	mode := uint32(0100644)
	if info.Mode()&0111 != 0 {
		// Executable bit set
		mode = uint32(0100755)
	}

	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		// Fallback if we can't get the stat information
		entry := IndexEntry{
			Ctime: time.Now(),
			Mtime: info.ModTime(),
			Dev:   0, // Default values
			Ino:   0, // Default values
			Mode:  mode,
			Uid:   0, // Default values
			Gid:   0, // Default values
			Size:  uint32(len(data)),
			Hash:  hash,
			Flags: uint16(len(path) & 0xFFF), // Path length in lower 12 bits
			Path:  path,
		}
		idx.Entries = append(idx.Entries, entry)
	} else {
		// Use explicit type conversion for the stat fields
		entry := IndexEntry{
			Ctime: time.Now(),
			Mtime: info.ModTime(),
			Dev:   uint32(stat.Dev), // Convert to uint32
			Ino:   uint32(stat.Ino), // Convert to uint32
			Mode:  mode,
			Uid:   uint32(stat.Uid), // Convert to uint32
			Gid:   uint32(stat.Gid), // Convert to uint32
			Size:  uint32(len(data)),
			Hash:  hash,
			Flags: uint16(len(path) & 0xFFF), // Path length in lower 12 bits
			Path:  path,
		}
		idx.Entries = append(idx.Entries, entry)
	}

	// Save index
	return idx.Save()
}

// Remove a file from the index
func (idx *Index) RemoveFile(path string) error {
	for i, entry := range idx.Entries {
		if entry.Path == path {
			// Remove entry
			idx.Entries = append(idx.Entries[:i], idx.Entries[i+1:]...)
			return idx.Save()
		}
	}
	return fmt.Errorf("file not in index: %s", path)
}

// Get staged files as a map of path to hash
func (idx *Index) GetStagedFiles() map[string]string {
	result := make(map[string]string)
	for _, entry := range idx.Entries {
		result[entry.Path] = entry.Hash
	}
	return result
}

// Check if a file is modified compared to the index
func (idx *Index) IsModified(path string) (bool, error) {
	// Find entry in index
	var entry *IndexEntry
	for i := range idx.Entries {
		if idx.Entries[i].Path == path {
			entry = &idx.Entries[i]
			break
		}
	}

	if entry == nil {
		// Not in index
		return true, nil
	}

	// Read current file content
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	// Hash current content
	currentHash, err := createBlob(data)
	if err != nil {
		return false, err
	}

	// Compare hash with index
	return currentHash != entry.Hash, nil
}

// Create a tree object from current index
func (idx *Index) WriteTree() (string, error) {
	// Group entries by directory
	directories := make(map[string][]IndexEntry)

	for _, entry := range idx.Entries {
		dir := filepath.Dir(entry.Path)
		if dir == "." {
			dir = ""
		}
		directories[dir] = append(directories[dir], entry)
	}

	// Build trees from bottom up
	return idx.buildTree("", directories)
}

// Recursively build tree objects
func (idx *Index) buildTree(prefix string, directories map[string][]IndexEntry) (string, error) {
	var entries [][]byte

	// Process direct children of this directory
	directEntries, ok := directories[prefix]
	if ok {
		for _, entry := range directEntries {
			// Get file name (basename)
			name := filepath.Base(entry.Path)

			// Determine mode
			mode := fmt.Sprintf("%o", entry.Mode)

			// Encode tree entry
			treeEntry, err := encodeTreeEntry(mode, name, entry.Hash)
			if err != nil {
				return "", err
			}

			entries = append(entries, treeEntry)
		}
	}

	// Process subdirectories
	for dir := range directories {
		// Check if dir is a direct subdirectory of prefix
		if dir != prefix && (prefix == "" || strings.HasPrefix(dir, prefix+"/")) {
			parts := strings.Split(dir, "/")
			dirName := parts[len(parts)-1]

			// Skip if already processed
			alreadyProcessed := false
			for _, entry := range entries {
				entryName := strings.SplitN(string(entry), "\x00", 2)[0]
				entryName = strings.SplitN(entryName, " ", 2)[1]
				if entryName == dirName {
					alreadyProcessed = true
					break
				}
			}

			if !alreadyProcessed {
				// Recursive build subdirectory tree
				subTreeHash, err := idx.buildTree(dir, directories)
				if err != nil {
					return "", err
				}

				// Add subdirectory entry to this tree
				treeEntry, err := encodeTreeEntry(DirectoryMode, dirName, subTreeHash)
				if err != nil {
					return "", err
				}

				entries = append(entries, treeEntry)
			}
		}
	}

	// Sort entries as Git does
	sort.Slice(entries, func(i, j int) bool {
		return bytes.Compare(entries[i], entries[j]) < 0
	})

	// Concatenate all entries
	var treeData []byte
	for _, entry := range entries {
		treeData = append(treeData, entry...)
	}

	// Create tree object
	return createTree(treeData)
}
