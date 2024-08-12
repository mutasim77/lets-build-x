package command

import (
	"redis/resp"
	"redis/storage"
)

// 1) -> https://redis.io/docs/latest/commands/ping
// Ping handles the PING command
// It returns "PONG" if no argument is provided, otherwise it returns the first argument
func Ping(s *storage.Storage, args []string) resp.Value {
	if len(args) == 0 {
		return resp.Value{Type: "string", Str: "PONG"}
	}
	return resp.Value{Type: "string", Str: args[0]}
}

// 2) -> https://redis.io/docs/latest/commands/set
// Set handles the SET command
// It sets a key to hold a string value in the storage
func Set(s *storage.Storage, args []string) resp.Value {
	if len(args) != 2 {
		return resp.Value{Type: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}
	s.Set(args[0], args[1])
	return resp.Value{Type: "string", Str: "OK"}
}

// 3) -> https://redis.io/docs/latest/commands/get
// Get handles the GET command
// It retrieves the value of a key from the storage
func Get(s *storage.Storage, args []string) resp.Value {
	if len(args) != 1 {
		return resp.Value{Type: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}
	value, ok := s.Get(args[0])
	if !ok {
		return resp.Value{Type: "null"}
	}
	return resp.Value{Type: "bulk", Bulk: value}
}

// 4) -> https://redis.io/docs/latest/commands/del
// Del handles the DEL command
// It removes the specified keys from the storage
// Returns the number of keys that were removed
func Del(s *storage.Storage, args []string) resp.Value {
	if len(args) < 1 {
		return resp.Value{Type: "error", Str: "ERR wrong number of arguments for 'del' command"}
	}
	count := 0
	for _, key := range args {
		if s.Del(key) {
			count++
		}
	}
	return resp.Value{Type: "integer", Num: count}
}

// 5) -> https://redis.io/docs/latest/commands/exists
// Exists handles the EXISTS command
// It checks if the specified keys exist in the storage
// Returns the number of keys that exist
func Exists(s *storage.Storage, args []string) resp.Value {
	if len(args) < 1 {
		return resp.Value{Type: "error", Str: "ERR wrong number of arguments for 'exists' command"}
	}
	count := s.Exists(args...)
	return resp.Value{Type: "integer", Num: count}
}

// 6) -> https://redis.io/docs/latest/commands/incr
// Incr handles the INCR command
// It increments the integer value of a key by one
// If the key does not exist, it is set to 0 before performing the operation
func Incr(s *storage.Storage, args []string) resp.Value {
    if len(args) != 1 {
        return resp.Value{Type: "error", Str: "ERR wrong number of arguments for 'incr' command"}
    }
    newValue, err := s.IncrBy(args[0], 1)
    if err != nil {
        return resp.Value{Type: "error", Str: err.Error()}
    }
    return resp.Value{Type: "integer", Num: newValue}
}
