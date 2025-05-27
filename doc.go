// Package sshserver provides a flexible and configurable SSH server implementation.
//
// This package allows you to create an SSH server with custom authentication,
// command handling, and logging capabilities. It's designed to be easy to use
// while remaining highly configurable.
//
// Basic usage:
//
//	config := sshserver.DefaultConfig()
//	config.ListenAddress = ":2222"
//	config.HostKeyFile = "path/to/host_key"
//	config.AuthorizedKeysFile = "path/to/authorized_keys"
//
//	handler := sshserver.NewDefaultHandler()
//	
//	server, err := sshserver.NewServer(config, handler)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if err := server.Start(); err != nil {
//		log.Fatal(err)
//	}
//
// Custom Command Handler:
//
//	type MyHandler struct{}
//	
//	func (h *MyHandler) Execute(cmd string) (string, uint32) {
//		// Handle command execution
//		return "Command output", 0
//	}
//	
//	func (h *MyHandler) GetPrompt() string {
//		return "my-server> "
//	}
//	
//	func (h *MyHandler) GetWelcomeMessage() string {
//		return "Welcome to My SSH Server!"
//	}
//
// Features:
//   - Public key authentication
//   - Custom command handling
//   - Configurable logging
//   - Graceful shutdown
//   - Interactive shell support
//   - Command execution support
//
// The package follows Go idioms and best practices, making it easy to integrate
// into existing projects while maintaining flexibility for custom implementations.
package sshserver