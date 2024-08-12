package server

import (
	"fmt"
	"net"
	"redis/aof"
	"redis/command"
	"redis/resp"
	"redis/storage"
)

// Server represents the Redis-like server
type Server struct {
	Addr    string           // Address to listen on
	Storage *storage.Storage // In-memory storage
	AOF     *aof.AOF         // Append-Only File for persistence
}

// NewServer creates a new Server instance
func NewServer(addr string) (*Server, error) {
	storage := storage.NewStorage()
	aofHandler, err := aof.NewAOF("database.aof")
	if err != nil {
		return nil, fmt.Errorf("failed to create AOF handler: %v", err)
	}

	server := &Server{
		Addr:    addr,
		Storage: storage,
		AOF:     aofHandler,
	}

	if err := server.loadAOF(); err != nil {
		return nil, fmt.Errorf("failed to load AOF: %v", err)
	}

	return server, nil
}

// loadAOF loads the Append-Only File and executes all commands
func (s *Server) loadAOF() error {
	return s.AOF.Load(func(value resp.Value) {
		if value.Type == "array" && len(value.Array) > 0 {
			cmd := value.Array[0].Bulk
			args := make([]string, len(value.Array)-1)
			for i, v := range value.Array[1:] {
				args[i] = v.Bulk
			}
			s.executeCommand(cmd, args)
		}
	})
}

// Run starts the server and listens for connections
func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Server listening on %s\n", s.Addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// handleConnection processes client connections
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	respReader := resp.NewResp(conn)

	for {
		value, err := respReader.Read()
		if err != nil {
			fmt.Printf("Error reading command: %v\n", err)
			return
		}

		if value.Type != "array" {
			fmt.Println("Invalid command format")
			continue
		}

		if len(value.Array) == 0 {
			fmt.Println("Empty command")
			continue
		}

		// Write command to AOF for persistence
		if err := s.AOF.Write(value); err != nil {
			fmt.Printf("Error writing to AOF: %v\n", err)
		}

		cmd := value.Array[0].Bulk
		args := make([]string, len(value.Array)-1)
		for i, v := range value.Array[1:] {
			args[i] = v.Bulk
		}

		result := s.executeCommand(cmd, args)
		if _, err := conn.Write(result.Marshal()); err != nil {
			fmt.Printf("Error writing response: %v\n", err)
			return
		}
	}
}

// executeCommand executes the given command with its arguments
func (s *Server) executeCommand(cmd string, args []string) resp.Value {
	switch cmd {
	case "PING":
		return command.Ping(s.Storage, args)
	case "SET":
		return command.Set(s.Storage, args)
	case "GET":
		return command.Get(s.Storage, args)
	case "DEL":
		return command.Del(s.Storage, args)
	case "EXISTS":
		return command.Exists(s.Storage, args)
	case "INCR":
		return command.Incr(s.Storage, args)
	default:
		return resp.Value{Type: "error", Str: "ERR unknown command '" + cmd + "'"}
	}
}
