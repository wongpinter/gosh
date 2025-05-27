package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"repo.nusatek.id/sugeng/gosh"
)

func main() {
	// Create default configuration
	config := sshserver.DefaultConfig()
	config.ListenAddress = ":2222"
	config.HostKeyFile = "server_key"
	config.AuthorizedKeysFile = "authorized_keys"
	
	// Enable logging to both file and stdout
	config.LogWriter.Enabled = true
	config.LogWriter.FilePath = "ssh_server.log"
	config.LogWriter.LogToStdout = true

	// Create default command handler
	handler := sshserver.NewDefaultHandler()

	// Create and start the server
	server, err := sshserver.NewServer(config, handler)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start the server
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("SSH server started successfully!")
	log.Println("Connect with: ssh -p 2222 user@localhost")
	log.Println("Available commands: hello, getDate, uptime, help")

	// Wait for interrupt signal to gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down server...")
	if err := server.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	}
	log.Println("Server stopped")
}
