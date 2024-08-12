package main

import (
	"fmt"
	"log"
	"redis/server"
)

func main() {
	// Define the port on which the server will listen
	const port = ":6379"

	// Print a startup message
	fmt.Printf("Starting Redis server on port %s\n", port)

	// Create a new server instance
	server, err := server.NewServer(port)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Run the server
	fmt.Println("Server is ready to accept connections")
	if err := server.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
