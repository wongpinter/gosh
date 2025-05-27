# SSH Server Package (gosh)

A flexible and configurable SSH server implementation in Go that allows you to create custom SSH servers with authentication, command handling, and logging capabilities.

## Features

* 🔐 **Public Key Authentication** - Secure SSH key-based authentication
* 🎯 **Custom Command Handlers** - Implement your own command processing logic
* 📝 **Configurable Logging** - Log to files, stdout, or both
* 🔄 **Graceful Shutdown** - Clean server termination with signal handling
* 🖥️ **Interactive Shell Support** - Full shell-like experience with prompts
* ⚡ **Command Execution** - Direct command execution without shell
* 🛡️ **Security First** - Built with security best practices
* 🔧 **Easy Configuration** - Simple configuration with sensible defaults

## Installation

```bash
go get repo.nusatek.id/sugeng/gosh
```

## Quick Start

```go
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

    // Create default command handler
    handler := sshserver.NewDefaultHandler()

    // Create and start server
    server, err := sshserver.NewServer(config, handler)
    if err != nil {
        log.Fatal(err)
    }

    if err := server.Start(); err != nil {
        log.Fatal(err)
    }

    // Wait for interrupt signal
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    <-c

    server.Stop()
}
```

## Configuration

### Basic Configuration

```go
config := sshserver.DefaultConfig()
config.ListenAddress = ":2222"              // Port to listen on
config.HostKeyFile = "server_key"           // SSH host private key
config.AuthorizedKeysFile = "authorized_keys" // Authorized public keys
```

### Advanced Configuration

```go
config := &sshserver.Config{
    ListenAddress:      ":2222",
    HostKeyFile:        "server_key",
    AuthorizedKeysFile: "authorized_keys",
    NoClientAuth:       false,  // Require authentication
    AllowKeyboardInteractive: false,
    LogWriter: &sshserver.LogConfig{
        Enabled:     true,
        FilePath:    "ssh_server.log",
        LogToStdout: true,
    },
}
```

## Custom Command Handlers

Implement the `CommandHandler` interface to create custom functionality:

```go
type MyHandler struct{}

func (h *MyHandler) Execute(cmd string) (string, uint32) {
    switch cmd {
    case "hello":
        return "Hello, World!", 0
    case "time":
        return time.Now().String(), 0
    default:
        return "Unknown command", 1
    }
}

func (h *MyHandler) GetPrompt() string {
    return "my-server> "
}

func (h *MyHandler) GetWelcomeMessage() string {
    return "Welcome to My Custom SSH Server!"
}
```

## Examples

The package includes comprehensive examples demonstrating various use cases:

### 🔧 Basic Server

Simple SSH server with default commands.

```bash
cd examples/basic-server
./setup.sh && go run main.go
ssh -p 2222 user@localhost
```

### ⚙️ Custom Handler

Advanced command processing with calculator and text tools.

```bash
cd examples/custom-handler
./setup.sh && go run main.go
ssh -p 2223 user@localhost
```

### 📂 File Server

Browse and manage files over SSH.

```bash
cd examples/file-server
./setup.sh && go run main.go
ssh -p 2224 user@localhost
```

### 🛠️ Admin Panel

System administration and monitoring.

```bash
cd examples/admin-panel
./setup.sh && go run main.go
ssh -p 2225 admin@localhost
```

### 💬 Chat Server

Multi-user real-time chat system.

```bash
cd examples/chat-server
./setup.sh && go run main.go
ssh -p 2226 alice@localhost  # Terminal 1
ssh -p 2226 bob@localhost    # Terminal 2
```

### 🎮 Game Server

Interactive games over SSH.

```bash
cd examples/game-server
./setup.sh && go run main.go
ssh -p 2227 player@localhost
```

### 📊 Monitoring Server

Real-time metrics and system monitoring.

```bash
cd examples/monitoring-server
./setup.sh && go run main.go
ssh -p 2228 monitor@localhost
```

## API Reference

### Core Types

#### Config

```go
type Config struct {
    ListenAddress      string     // Address to listen on (e.g., ":2222")
    HostKeyFile        string     // Path to SSH host private key
    AuthorizedKeysFile string     // Path to authorized_keys file
    NoClientAuth       bool       // Disable client authentication
    AllowKeyboardInteractive bool // Enable keyboard-interactive auth
    LogWriter          *LogConfig // Logging configuration
}
```

#### CommandHandler Interface

```go
type CommandHandler interface {
    Execute(cmd string) (string, uint32)  // Process command, return output and exit code
    GetPrompt() string                    // Return shell prompt
    GetWelcomeMessage() string           // Return welcome message
}
```

#### Server

```go
type Server struct {
    // Server manages SSH connections and command processing
}

func NewServer(config *Config, handler CommandHandler) (*Server, error)
func (s *Server) Start() error
func (s *Server) Stop() error
```

### Default Handler

The package provides a default command handler with basic commands:

```go
handler := sshserver.NewDefaultHandler()
```

**Built-in Commands:**
* `hello` - Simple greeting
* `getDate` - Current server time
* `uptime` - Server uptime information
* `help` - List available commands

### Custom Commands

Add custom commands to the default handler:

```go
handler := sshserver.NewDefaultHandler()
handler.RegisterCommand("custom", func() (string, error) {
    return "Custom command output", nil
})
```

## Security

### SSH Key Generation

Generate server host key:

```bash
ssh-keygen -t rsa -b 2048 -f server_key -N ""
```

Generate client key:

```bash
ssh-keygen -t rsa -b 2048 -f client_key -N ""
```

Add client public key to authorized_keys:

```bash
cp client_key.pub authorized_keys
```

### Best Practices

1. **Use Strong Keys** - Generate 2048-bit or larger RSA keys
2. **Secure Key Storage** - Protect private keys with proper permissions (600)
3. **Regular Key Rotation** - Rotate SSH keys periodically
4. **Audit Logging** - Enable comprehensive logging for security auditing
5. **Network Security** - Use firewalls and network segmentation
6. **Input Validation** - Validate all command inputs in custom handlers

## Use Cases

### Development & Testing

* Mock SSH services for testing
* Development environment simulation
* API testing over SSH

### System Administration

* Remote server management
* Automated deployment tools
* System monitoring interfaces

### Interactive Applications

* Chat systems and communication tools
* Interactive games and entertainment
* Educational and training platforms

### File Management

* Secure file transfer systems
* Remote file browsing
* Backup and synchronization tools

### Monitoring & Analytics

* Real-time system monitoring
* Metrics collection and reporting
* Alert and notification systems

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For questions, issues, or contributions:
* Create an issue on GitHub
* Check the examples directory for usage patterns
* Review the documentation and API reference

## Example Showcase

### File Server Example

Create a secure file browser accessible over SSH:

```go
// Browse files
files:/> ls
Directory: /
d      <DIR> 2024-01-15 10:30 documents
-      85 B 2024-01-15 10:30 config.json

// Navigate directories
files:/> cd documents
Changed directory to: /documents

// View file contents
files:/documents> cat notes.txt
Meeting Notes
=============
1. Project status
2. Next steps

// Search for files
files:/> find *.txt
Found 2 matches for pattern '*.txt':
  /readme.txt
  /documents/notes.txt
```

### Chat Server Example

Multi-user real-time chat over SSH:

```bash
# Terminal 1 (Alice)
[alice] Hello everyone!
[15:30:22] * bob joined the chat
[alice] /users
Online users (2):
  alice (you)
  bob

# Terminal 2 (Bob)
[bob] Hey Alice! How's it going?
[bob] /me waves
[15:30:45] * bob waves
```

### Game Server Example

Interactive games accessible via SSH:

```bash
🎮 game> guess
🔢 NUMBER GUESSING GAME 🔢
Make your first guess!

🔢 guess> 50
📉 Too high! You have 6 attempts left.

🔢 guess> 25
📈 Too low! You have 5 attempts left.

🔢 guess> 33
🎉 Congratulations! You earned 20 points!
```

### Monitoring Server Example

Real-time system monitoring:

```bash
monitor> dashboard
=== MONITORING DASHBOARD ===
Uptime: 5m30s
Memory: 2.1 MB / 8.7 MB
Goroutines: 12
Requests: 15

monitor> alert
✅ All systems normal - no alerts

monitor> export json
=== METRICS EXPORT (JSON) ===
[
  {
    "timestamp": "2024-01-15T15:30:15Z",
    "type": "memory.alloc",
    "value": 2097152,
    "unit": "bytes"
  }
]
```

## Advanced Usage

### Multiple Servers

Run multiple SSH servers on different ports:

```go
// Create multiple configurations
configs := []*sshserver.Config{
    {ListenAddress: ":2222", HostKeyFile: "key1"},
    {ListenAddress: ":2223", HostKeyFile: "key2"},
}

// Start multiple servers
for _, config := range configs {
    server, _ := sshserver.NewServer(config, handler)
    go server.Start()
}
```

### Dynamic Command Registration

Add commands at runtime:

```go
handler := sshserver.NewDefaultHandler()

// Add commands dynamically
handler.RegisterCommand("status", func() (string, error) {
    return "Server is running", nil
})

handler.RegisterCommand("users", func() (string, error) {
    return fmt.Sprintf("Active users: %d", getUserCount()), nil
})
```

### Session Management

Track user sessions and state:

```go
type SessionHandler struct {
    sessions map[string]*UserSession
    mutex    sync.RWMutex
}

func (h *SessionHandler) Execute(cmd string) (string, uint32) {
    // Access session-specific data
    session := h.getCurrentSession()
    return h.processCommand(cmd, session)
}
```

### Integration Examples

#### With Web Servers

```go
// Embed SSH server in web application
go func() {
    sshServer, _ := sshserver.NewServer(sshConfig, handler)
    sshServer.Start()
}()

// Continue with HTTP server
http.ListenAndServe(":8080", webHandler)
```

#### With Databases

```go
type DBHandler struct {
    db *sql.DB
}

func (h *DBHandler) Execute(cmd string) (string, uint32) {
    // Execute database queries via SSH
    rows, err := h.db.Query("SELECT * FROM users")
    // Process and return results
}
```

#### With External APIs

```go
func (h *APIHandler) Execute(cmd string) (string, uint32) {
    // Proxy commands to external APIs
    resp, err := http.Get("https://api.example.com/data")
    // Return API response
}
```

## Performance Considerations

### Connection Limits

```go
config.MaxConnections = 100  // Limit concurrent connections
```

### Memory Management

```go
// Implement cleanup in handlers
func (h *Handler) cleanup() {
    // Clean up resources
    h.cache.Clear()
    runtime.GC()
}
```

### Logging Performance

```go
// Use buffered logging for high-throughput scenarios
config.LogWriter = &LogConfig{
    Enabled:     true,
    FilePath:    "server.log",
    BufferSize:  8192,  // Buffer log writes
}
```

## Testing

### Unit Tests

```go
func TestCustomHandler(t *testing.T) {
    handler := NewCustomHandler()
    output, code := handler.Execute("test")

    assert.Equal(t, "test output", output)
    assert.Equal(t, uint32(0), code)
}
```

### Integration Tests

```go
func TestSSHServer(t *testing.T) {
    // Start test server
    server, _ := sshserver.NewServer(testConfig, testHandler)
    go server.Start()
    defer server.Stop()

    // Test SSH connection
    client := ssh.Dial("tcp", "localhost:2222", clientConfig)
    session, _ := client.NewSession()

    output, _ := session.CombinedOutput("test command")
    assert.Contains(t, string(output), "expected")
}
```

## Troubleshooting

### Common Issues

**Connection Refused**

```bash
# Check if server is running
netstat -ln | grep 2222

# Verify host key exists
ls -la server_key
```

**Authentication Failed**

```bash
# Check authorized_keys format
ssh-keygen -l -f authorized_keys

# Verify key permissions
chmod 600 client_key
chmod 644 authorized_keys
```

**Command Not Found**

```go
// Ensure command is registered
handler.RegisterCommand("mycommand", handlerFunc)

// Check command parsing
fmt.Printf("Received command: %q\n", cmd)
```

### Debug Mode

```go
config.LogWriter.LogLevel = "DEBUG"  // Enable debug logging
config.LogWriter.LogToStdout = true  // See logs in console
```

## Changelog

### v1.0.0

* Initial release
* Basic SSH server functionality
* Public key authentication
* Custom command handlers
* Comprehensive examples
* Full documentation
