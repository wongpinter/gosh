package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	sshserver "repo.nusatek.id/sugeng/gosh"
)

// AdminHandler implements administrative commands
type AdminHandler struct {
	startTime    time.Time
	commandCount int
	allowedUsers map[string]bool
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		startTime:    time.Now(),
		commandCount: 0,
		allowedUsers: map[string]bool{
			"admin": true,
			"root":  true,
		},
	}
}

// Execute implements the CommandHandler interface
func (h *AdminHandler) Execute(cmd string) (string, uint32) {
	h.commandCount++
	parts := strings.Fields(strings.TrimSpace(cmd))
	if len(parts) == 0 {
		return "", 0
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "status":
		return h.getSystemStatus(), 0
	case "uptime":
		return h.getUptime(), 0
	case "memory", "mem":
		return h.getMemoryInfo(), 0
	case "disk":
		return h.getDiskInfo(), 0
	case "processes", "ps":
		return h.getProcesses(args)
	case "network", "net":
		return h.getNetworkInfo(), 0
	case "users":
		return h.getLoggedUsers(), 0
	case "logs":
		return h.getSystemLogs(args)
	case "services":
		return h.getServices(args)
	case "load":
		return h.getLoadAverage(), 0
	case "env":
		return h.getEnvironment(args)
	case "date":
		return h.getDateTime(), 0
	case "whoami":
		return h.getCurrentUser(), 0
	case "stats":
		return h.getServerStats(), 0
	case "help":
		return h.getHelp(), 0
	default:
		return fmt.Sprintf("Unknown command: %s\nType 'help' for available commands", command), 1
	}
}

func (h *AdminHandler) getSystemStatus() string {
	var result strings.Builder
	result.WriteString("=== SYSTEM STATUS ===\n")
	result.WriteString(fmt.Sprintf("OS: %s\n", runtime.GOOS))
	result.WriteString(fmt.Sprintf("Architecture: %s\n", runtime.GOARCH))
	result.WriteString(fmt.Sprintf("Go Version: %s\n", runtime.Version()))
	result.WriteString(fmt.Sprintf("CPUs: %d\n", runtime.NumCPU()))
	result.WriteString(fmt.Sprintf("Goroutines: %d\n", runtime.NumGoroutine()))

	// Get hostname
	if hostname, err := os.Hostname(); err == nil {
		result.WriteString(fmt.Sprintf("Hostname: %s\n", hostname))
	}

	return result.String()
}

func (h *AdminHandler) getUptime() string {
	if runtime.GOOS == "linux" {
		if output, err := exec.Command("uptime").Output(); err == nil {
			return strings.TrimSpace(string(output))
		}
	}

	// Fallback: show server uptime
	uptime := time.Since(h.startTime)
	return fmt.Sprintf("Server uptime: %v", uptime.Round(time.Second))
}

func (h *AdminHandler) getMemoryInfo() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	var result strings.Builder
	result.WriteString("=== MEMORY INFO ===\n")
	result.WriteString(fmt.Sprintf("Allocated: %s\n", h.formatBytes(m.Alloc)))
	result.WriteString(fmt.Sprintf("Total Allocated: %s\n", h.formatBytes(m.TotalAlloc)))
	result.WriteString(fmt.Sprintf("System: %s\n", h.formatBytes(m.Sys)))
	result.WriteString(fmt.Sprintf("GC Runs: %d\n", m.NumGC))

	// Try to get system memory info on Linux
	if runtime.GOOS == "linux" {
		if output, err := exec.Command("free", "-h").Output(); err == nil {
			result.WriteString("\n=== SYSTEM MEMORY ===\n")
			result.WriteString(string(output))
		}
	}

	return result.String()
}

func (h *AdminHandler) getDiskInfo() string {
	if runtime.GOOS == "linux" {
		if output, err := exec.Command("df", "-h").Output(); err == nil {
			return "=== DISK USAGE ===\n" + string(output)
		}
	}

	// Fallback: show current directory info
	if pwd, err := os.Getwd(); err == nil {
		if stat, err := os.Stat(pwd); err == nil {
			return fmt.Sprintf("Current directory: %s\nLast modified: %s",
				pwd, stat.ModTime().Format("2006-01-02 15:04:05"))
		}
	}

	return "Disk information not available on this platform"
}

func (h *AdminHandler) getProcesses(args []string) (string, uint32) {
	if runtime.GOOS != "linux" {
		return "Process listing not available on this platform", 1
	}

	var cmd *exec.Cmd
	if len(args) > 0 && args[0] == "full" {
		cmd = exec.Command("ps", "aux")
	} else {
		cmd = exec.Command("ps", "-eo", "pid,ppid,user,comm,%cpu,%mem")
	}

	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error getting processes: %v", err), 1
	}

	return "=== PROCESSES ===\n" + string(output), 0
}

func (h *AdminHandler) getNetworkInfo() string {
	var result strings.Builder
	result.WriteString("=== NETWORK INFO ===\n")

	if runtime.GOOS == "linux" {
		// Get network interfaces
		if output, err := exec.Command("ip", "addr", "show").Output(); err == nil {
			result.WriteString("Network Interfaces:\n")
			result.WriteString(string(output))
			result.WriteString("\n")
		}

		// Get network connections
		if output, err := exec.Command("ss", "-tuln").Output(); err == nil {
			result.WriteString("Listening Ports:\n")
			result.WriteString(string(output))
		}
	} else {
		result.WriteString("Network information not available on this platform")
	}

	return result.String()
}

func (h *AdminHandler) getLoggedUsers() string {
	if runtime.GOOS == "linux" {
		if output, err := exec.Command("who").Output(); err == nil {
			return "=== LOGGED USERS ===\n" + string(output)
		}
	}

	return "User information not available on this platform"
}

func (h *AdminHandler) getSystemLogs(args []string) (string, uint32) {
	if runtime.GOOS != "linux" {
		return "System logs not available on this platform", 1
	}

	lines := "20"
	if len(args) > 0 {
		if _, err := strconv.Atoi(args[0]); err == nil {
			lines = args[0]
		}
	}

	cmd := exec.Command("journalctl", "-n", lines, "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		// Fallback to syslog
		cmd = exec.Command("tail", "-n", lines, "/var/log/syslog")
		if output, err = cmd.Output(); err != nil {
			return fmt.Sprintf("Error reading logs: %v", err), 1
		}
	}

	return fmt.Sprintf("=== SYSTEM LOGS (last %s lines) ===\n%s", lines, string(output)), 0
}

func (h *AdminHandler) getServices(args []string) (string, uint32) {
	if runtime.GOOS != "linux" {
		return "Service information not available on this platform", 1
	}

	var cmd *exec.Cmd
	if len(args) > 0 {
		cmd = exec.Command("systemctl", "status", args[0])
	} else {
		cmd = exec.Command("systemctl", "list-units", "--type=service", "--state=running")
	}

	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error getting service info: %v", err), 1
	}

	return "=== SERVICES ===\n" + string(output), 0
}

func (h *AdminHandler) getLoadAverage() string {
	if runtime.GOOS == "linux" {
		if data, err := os.ReadFile("/proc/loadavg"); err == nil {
			return "=== LOAD AVERAGE ===\n" + string(data)
		}
	}

	return "Load average not available on this platform"
}

func (h *AdminHandler) getEnvironment(args []string) (string, uint32) {
	if len(args) > 0 {
		// Show specific environment variable
		value := os.Getenv(args[0])
		if value == "" {
			return fmt.Sprintf("Environment variable '%s' not set", args[0]), 0
		}
		return fmt.Sprintf("%s=%s", args[0], value), 0
	}

	// Show all environment variables
	var result strings.Builder
	result.WriteString("=== ENVIRONMENT VARIABLES ===\n")

	envVars := os.Environ()
	for _, env := range envVars {
		result.WriteString(env + "\n")
	}

	return result.String(), 0
}

func (h *AdminHandler) getDateTime() string {
	now := time.Now()
	return fmt.Sprintf("Current date/time: %s\nUnix timestamp: %d\nTimezone: %s",
		now.Format("2006-01-02 15:04:05 MST"),
		now.Unix(),
		now.Location().String())
}

func (h *AdminHandler) getCurrentUser() string {
	if user := os.Getenv("USER"); user != "" {
		return fmt.Sprintf("Current user: %s", user)
	}
	if user := os.Getenv("USERNAME"); user != "" {
		return fmt.Sprintf("Current user: %s", user)
	}
	return "Current user: unknown"
}

func (h *AdminHandler) getServerStats() string {
	uptime := time.Since(h.startTime)
	return fmt.Sprintf("=== SERVER STATISTICS ===\n"+
		"Server uptime: %v\n"+
		"Commands executed: %d\n"+
		"Started at: %s\n"+
		"Go version: %s\n"+
		"Platform: %s/%s",
		uptime.Round(time.Second),
		h.commandCount,
		h.startTime.Format("2006-01-02 15:04:05"),
		runtime.Version(),
		runtime.GOOS, runtime.GOARCH)
}

func (h *AdminHandler) formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func (h *AdminHandler) getHelp() string {
	return `Administrative Commands:
- status                 Show system status
- uptime                 Show system uptime
- memory, mem            Show memory information
- disk                   Show disk usage
- processes, ps [full]   Show running processes
- network, net           Show network information
- users                  Show logged in users
- logs [lines]           Show system logs (default: 20 lines)
- services [name]        Show services status
- load                   Show load average
- env [variable]         Show environment variables
- date                   Show current date/time
- whoami                 Show current user
- stats                  Show server statistics
- help                   Show this help message

Note: Some commands are platform-specific and may not work on all systems.`
}

// GetPrompt implements the CommandHandler interface
func (h *AdminHandler) GetPrompt() string {
	return "admin# "
}

// GetWelcomeMessage implements the CommandHandler interface
func (h *AdminHandler) GetWelcomeMessage() string {
	hostname, _ := os.Hostname()
	return fmt.Sprintf("Welcome to Admin Panel SSH Server!\n"+
		"Hostname: %s\n"+
		"Platform: %s/%s\n"+
		"Type 'help' to see available administrative commands.\n"+
		"Type 'status' for a quick system overview.",
		hostname, runtime.GOOS, runtime.GOARCH)
}

func main() {
	// Create configuration
	config := sshserver.DefaultConfig()
	config.ListenAddress = ":2225"
	config.HostKeyFile = "server_key"
	config.AuthorizedKeysFile = "authorized_keys"
	config.LogWriter.FilePath = "admin_server.log"

	// Create admin handler
	handler := NewAdminHandler()

	// Create and start server
	server, err := sshserver.NewServer(config, handler)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Admin Panel SSH server started on port 2225!")
	log.Println("Connect with: ssh -p 2225 admin@localhost")

	// Wait for interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down server...")
	if err := server.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	}
}
