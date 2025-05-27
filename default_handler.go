package sshserver

import (
	"fmt"
	"strings"
	"syscall"
	"time"
)

// DefaultCommandHandler provides a basic implementation of CommandHandler
type DefaultCommandHandler struct {
	commands map[string]func() (string, error)
}

// NewDefaultHandler creates a new DefaultCommandHandler with basic commands
func NewDefaultHandler() *DefaultCommandHandler {
	h := &DefaultCommandHandler{
		commands: make(map[string]func() (string, error)),
	}

	// Register default commands
	h.RegisterCommand("hello", func() (string, error) {
		return "Hello from SSH Server!", nil
	})

	h.RegisterCommand("getDate", func() (string, error) {
		return fmt.Sprintf("Current server time: %s", time.Now().Format(time.RFC3339)), nil
	})

	h.RegisterCommand("uptime", func() (string, error) {
		var info syscall.Sysinfo_t
		err := syscall.Sysinfo(&info)
		if err != nil {
			return "", fmt.Errorf("error getting system info: %v", err)
		}

		uptime := time.Duration(info.Uptime) * time.Second
		days := int(uptime.Hours() / 24)
		hours := int(uptime.Hours()) % 24
		minutes := int(uptime.Minutes()) % 60

		return fmt.Sprintf("Server uptime: %d days, %d hours, %d minutes", days, hours, minutes), nil
	})

	h.RegisterCommand("help", func() (string, error) {
		var commands []string
		for cmd := range h.commands {
			commands = append(commands, cmd)
		}
		return fmt.Sprintf("Available commands: %s", strings.Join(commands, ", ")), nil
	})

	return h
}

// RegisterCommand adds a new command to the handler
func (h *DefaultCommandHandler) RegisterCommand(name string, handler func() (string, error)) {
	h.commands[name] = handler
}

// Execute implements CommandHandler.Execute
func (h *DefaultCommandHandler) Execute(cmd string) (string, uint32) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return "", 0
	}

	if handler, ok := h.commands[cmd]; ok {
		output, err := handler()
		if err != nil {
			return fmt.Sprintf("Error: %v", err), 1
		}
		return output, 0
	}

	return fmt.Sprintf("Unknown command: %s\nUse 'help' to see available commands", cmd), 1
}

// GetPrompt implements CommandHandler.GetPrompt
func (h *DefaultCommandHandler) GetPrompt() string {
	return "$ "
}

// GetWelcomeMessage implements CommandHandler.GetWelcomeMessage
func (h *DefaultCommandHandler) GetWelcomeMessage() string {
	return "Welcome to SSH Server!\nType 'help' to see available commands"
}