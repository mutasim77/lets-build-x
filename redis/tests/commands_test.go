package tests

import (
	"redis/command"
	"redis/storage"
	"testing"
)

// TestPing tests the PING command
func TestPing(t *testing.T) {
	s := storage.NewStorage()

	// Test PING without argument
	result := command.Ping(s, []string{})
	if result.Type != "string" || result.Str != "PONG" {
		t.Errorf("Expected PONG, got %v", result)
	}

	// Test PING with argument
	result = command.Ping(s, []string{"Hello"})
	if result.Type != "string" || result.Str != "Hello" {
		t.Errorf("Expected Hello, got %v", result)
	}
}

// TestSetAndGet tests the SET and GET commands
func TestSetAndGet(t *testing.T) {
	s := storage.NewStorage()

	// Test SET
	setResult := command.Set(s, []string{"key", "value"})
	if setResult.Type != "string" || setResult.Str != "OK" {
		t.Errorf("SET: Expected OK, got %v", setResult)
	}

	// Test GET
	getResult := command.Get(s, []string{"key"})
	if getResult.Type != "bulk" || getResult.Bulk != "value" {
		t.Errorf("GET: Expected value, got %v", getResult)
	}

	// Test GET non-existent key
	getResult = command.Get(s, []string{"nonexistent"})
	if getResult.Type != "null" {
		t.Errorf("GET nonexistent: Expected null, got %v", getResult)
	}
}

// TestDel tests the DEL command
func TestDel(t *testing.T) {
	s := storage.NewStorage()

	command.Set(s, []string{"key1", "value1"})
	command.Set(s, []string{"key2", "value2"})

	// Test DEL existing keys
	delResult := command.Del(s, []string{"key1", "key2"})
	if delResult.Type != "integer" || delResult.Num != 2 {
		t.Errorf("DEL: Expected 2, got %v", delResult)
	}

	// Test DEL non-existent key
	delResult = command.Del(s, []string{"nonexistent"})
	if delResult.Type != "integer" || delResult.Num != 0 {
		t.Errorf("DEL nonexistent: Expected 0, got %v", delResult)
	}
}

// TestExists tests the EXISTS command
func TestExists(t *testing.T) {
	s := storage.NewStorage()

	command.Set(s, []string{"key1", "value1"})
	command.Set(s, []string{"key2", "value2"})

	// Test EXISTS on existing keys
	existsResult := command.Exists(s, []string{"key1", "key2"})
	if existsResult.Type != "integer" || existsResult.Num != 2 {
		t.Errorf("EXISTS: Expected 2, got %v", existsResult)
	}

	// Test EXISTS on mix of existing and non-existing keys
	existsResult = command.Exists(s, []string{"key1", "nonexistent"})
	if existsResult.Type != "integer" || existsResult.Num != 1 {
		t.Errorf("EXISTS mix: Expected 1, got %v", existsResult)
	}
}

// TestIncr tests the INCR command
func TestIncr(t *testing.T) {
	s := storage.NewStorage()

	// Test INCR on non-existent key
	incrResult := command.Incr(s, []string{"counter"})
	if incrResult.Type != "integer" || incrResult.Num != 1 {
		t.Errorf("INCR new: Expected 1, got %v", incrResult)
	}

	// Test INCR on existing key
	incrResult = command.Incr(s, []string{"counter"})
	if incrResult.Type != "integer" || incrResult.Num != 2 {
		t.Errorf("INCR existing: Expected 2, got %v", incrResult)
	}

	// Test INCR on non-integer value
	command.Set(s, []string{"string", "hello"})
	incrResult = command.Incr(s, []string{"string"})
	if incrResult.Type != "error" {
		t.Errorf("INCR non-integer: Expected error, got %v", incrResult)
	}
}
