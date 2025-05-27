package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"repo.nusatek.id/sugeng/gosh"
)

// CustomHandler implements a specialized command handler
type CustomHandler struct {
	startTime time.Time
	counter   int
}

// NewCustomHandler creates a new custom handler
func NewCustomHandler() *CustomHandler {
	return &CustomHandler{
		startTime: time.Now(),
		counter:   0,
	}
}

// Execute implements the CommandHandler interface
func (h *CustomHandler) Execute(cmd string) (string, uint32) {
	h.counter++
	parts := strings.Fields(strings.TrimSpace(cmd))
	if len(parts) == 0 {
		return "", 0
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "echo":
		if len(args) == 0 {
			return "Usage: echo <message>", 1
		}
		return strings.Join(args, " "), 0

	case "calc":
		return h.handleCalculator(args)

	case "random":
		return h.handleRandom(args)

	case "stats":
		return h.handleStats(), 0

	case "time":
		return h.handleTime(args)

	case "reverse":
		if len(args) == 0 {
			return "Usage: reverse <text>", 1
		}
		text := strings.Join(args, " ")
		return h.reverseString(text), 0

	case "upper":
		if len(args) == 0 {
			return "Usage: upper <text>", 1
		}
		return strings.ToUpper(strings.Join(args, " ")), 0

	case "lower":
		if len(args) == 0 {
			return "Usage: lower <text>", 1
		}
		return strings.ToLower(strings.Join(args, " ")), 0

	case "help":
		return h.getHelp(), 0

	default:
		return fmt.Sprintf("Unknown command: %s\nType 'help' for available commands", command), 1
	}
}

func (h *CustomHandler) handleCalculator(args []string) (string, uint32) {
	if len(args) != 3 {
		return "Usage: calc <number1> <operator> <number2>\nOperators: +, -, *, /", 1
	}

	num1, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return fmt.Sprintf("Invalid number: %s", args[0]), 1
	}

	num2, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return fmt.Sprintf("Invalid number: %s", args[2]), 1
	}

	operator := args[1]
	var result float64

	switch operator {
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
	case "*":
		result = num1 * num2
	case "/":
		if num2 == 0 {
			return "Error: Division by zero", 1
		}
		result = num1 / num2
	default:
		return fmt.Sprintf("Unknown operator: %s", operator), 1
	}

	return fmt.Sprintf("%.2f %s %.2f = %.2f", num1, operator, num2, result), 0
}

func (h *CustomHandler) handleRandom(args []string) (string, uint32) {
	if len(args) == 0 {
		// Random number between 1-100
		return fmt.Sprintf("Random number: %d", rand.Intn(100)+1), 0
	}

	if len(args) == 1 {
		max, err := strconv.Atoi(args[0])
		if err != nil {
			return "Usage: random [max] or random <min> <max>", 1
		}
		return fmt.Sprintf("Random number (1-%d): %d", max, rand.Intn(max)+1), 0
	}

	if len(args) == 2 {
		min, err1 := strconv.Atoi(args[0])
		max, err2 := strconv.Atoi(args[1])
		if err1 != nil || err2 != nil {
			return "Usage: random [max] or random <min> <max>", 1
		}
		if min >= max {
			return "Error: min must be less than max", 1
		}
		result := rand.Intn(max-min+1) + min
		return fmt.Sprintf("Random number (%d-%d): %d", min, max, result), 0
	}

	return "Usage: random [max] or random <min> <max>", 1
}

func (h *CustomHandler) handleStats() string {
	uptime := time.Since(h.startTime)
	return fmt.Sprintf("Server Statistics:\n"+
		"- Uptime: %v\n"+
		"- Commands executed: %d\n"+
		"- Started at: %s",
		uptime.Round(time.Second),
		h.counter,
		h.startTime.Format("2006-01-02 15:04:05"))
}

func (h *CustomHandler) handleTime(args []string) (string, uint32) {
	now := time.Now()
	
	if len(args) == 0 {
		return fmt.Sprintf("Current time: %s", now.Format("2006-01-02 15:04:05")), 0
	}

	format := strings.Join(args, " ")
	switch format {
	case "unix":
		return fmt.Sprintf("Unix timestamp: %d", now.Unix()), 0
	case "iso":
		return fmt.Sprintf("ISO format: %s", now.Format(time.RFC3339)), 0
	case "rfc":
		return fmt.Sprintf("RFC format: %s", now.Format(time.RFC822)), 0
	default:
		return fmt.Sprintf("Unknown time format: %s\nAvailable: unix, iso, rfc", format), 1
	}
}

func (h *CustomHandler) reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func (h *CustomHandler) getHelp() string {
	return `Available Commands:
- echo <message>          Echo back the message
- calc <n1> <op> <n2>     Calculator (+, -, *, /)
- random [max]            Generate random number
- random <min> <max>      Generate random number in range
- stats                   Show server statistics
- time [format]           Show current time (unix, iso, rfc)
- reverse <text>          Reverse the text
- upper <text>            Convert to uppercase
- lower <text>            Convert to lowercase
- help                    Show this help message`
}

// GetPrompt implements the CommandHandler interface
func (h *CustomHandler) GetPrompt() string {
	return "custom> "
}

// GetWelcomeMessage implements the CommandHandler interface
func (h *CustomHandler) GetWelcomeMessage() string {
	return "Welcome to Custom SSH Server!\n" +
		"This server has enhanced commands for calculations, text processing, and more.\n" +
		"Type 'help' to see all available commands."
}

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Create configuration
	config := sshserver.DefaultConfig()
	config.ListenAddress = ":2223"
	config.HostKeyFile = "server_key"
	config.AuthorizedKeysFile = "authorized_keys"
	config.LogWriter.FilePath = "custom_server.log"

	// Create custom handler
	handler := NewCustomHandler()

	// Create and start server
	server, err := sshserver.NewServer(config, handler)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Custom SSH server started on port 2223!")
	log.Println("Connect with: ssh -p 2223 user@localhost")

	// Wait for interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down server...")
	if err := server.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	}
}
