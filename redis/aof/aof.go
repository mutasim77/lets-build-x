// https://redis.io/docs/latest/operate/oss_and_stack/management/persistence/
package aof

import (
	"bufio"
	"os"
	"redis/resp"
	"sync"
)

// AOF represents the Append-Only File structure for data persistence
type AOF struct {
	file   *os.File
	writer *bufio.Writer
	mu     sync.Mutex
}

// NewAOF creates a new AOF instance with the given filename
// It opens (or creates) the file and sets up a buffered writer
func NewAOF(filename string) (*AOF, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &AOF{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

// Write appends a RESP value to the AOF
// It uses a mutex to ensure thread-safety
func (aof *AOF) Write(value resp.Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	if _, err := aof.writer.Write(value.Marshal()); err != nil {
		return err
	}

	return aof.writer.Flush()
}

// Close flushes any remaining data and closes the AOF file
// It uses a mutex to ensure thread-safety
func (aof *AOF) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	if err := aof.writer.Flush(); err != nil {
		return err
	}

	return aof.file.Close()
}

// Load reads the AOF file from the beginning and applies each command
// using the provided handler function
// It uses a mutex to ensure thread-safety
func (aof *AOF) Load(handler func(resp.Value)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	if _, err := aof.file.Seek(0, 0); err != nil {
		return err
	}

	reader := resp.NewResp(aof.file)

	for {
		value, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		handler(value)
	}

	return nil
}
