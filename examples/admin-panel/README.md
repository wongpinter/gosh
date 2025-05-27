# Admin Panel SSH Server Example

This example demonstrates an SSH server that provides administrative and system monitoring capabilities, similar to a system administration panel accessible via SSH.

## Features

- System status monitoring
- Memory and disk usage information
- Process management and monitoring
- Network interface information
- System logs viewing
- Service status checking
- Load average monitoring
- Environment variable inspection
- User session information

## Setup

1. Generate SSH keys:
   ```bash
   ./setup.sh
   ```

2. Run the server:
   ```bash
   go run main.go
   ```

3. Connect as admin:
   ```bash
   ssh -p 2225 admin@localhost
   ```

## Available Commands

### System Information
- `status` - Show comprehensive system status
- `uptime` - Show system uptime
- `date` - Show current date/time with timezone
- `whoami` - Show current user
- `stats` - Show server statistics

### Resource Monitoring
- `memory`, `mem` - Show memory usage (Go runtime + system)
- `disk` - Show disk usage information
- `load` - Show system load average
- `processes`, `ps [full]` - Show running processes

### Network & Users
- `network`, `net` - Show network interfaces and listening ports
- `users` - Show currently logged in users

### System Management
- `logs [lines]` - Show system logs (default: 20 lines)
- `services [name]` - Show service status (all or specific service)
- `env [variable]` - Show environment variables (all or specific)

### Utility
- `help` - Show all available commands

## Example Session

```
$ ssh -p 2225 admin@localhost
Welcome to Admin Panel SSH Server!
Hostname: myserver
Platform: linux/amd64
Type 'help' to see available administrative commands.
Type 'status' for a quick system overview.

admin# status
=== SYSTEM STATUS ===
OS: linux
Architecture: amd64
Go Version: go1.24.3
CPUs: 8
Goroutines: 12
Hostname: myserver

admin# memory
=== MEMORY INFO ===
Allocated: 2.1 MB
Total Allocated: 3.4 MB
System: 8.7 MB
GC Runs: 5

=== SYSTEM MEMORY ===
              total        used        free      shared  buff/cache   available
Mem:           15Gi       2.1Gi        11Gi       234Mi       2.0Gi        13Gi
Swap:         2.0Gi          0B       2.0Gi

admin# processes
=== PROCESSES ===
    PID    PPID USER     COMMAND          %CPU %MEM
      1       0 root     systemd           0.0  0.1
    123       1 root     sshd              0.0  0.2
   1234     123 user     go                1.2  0.5

admin# disk
=== DISK USAGE ===
Filesystem      Size  Used Avail Use% Mounted on
/dev/sda1        20G  8.5G   11G  45% /
/dev/sda2       100G   45G   50G  48% /home

admin# logs 5
=== SYSTEM LOGS (last 5 lines) ===
Jan 15 10:30:45 myserver systemd[1]: Started SSH server.
Jan 15 10:31:02 myserver sshd[1234]: Accepted publickey for admin
Jan 15 10:31:15 myserver kernel: TCP: request_sock_TCP: Possible SYN flooding
```

## Platform Support

### Linux (Full Support)
- All commands are available
- Uses system utilities: `ps`, `df`, `free`, `who`, `journalctl`, `systemctl`
- Provides detailed system information

### Other Platforms (Limited Support)
- Basic Go runtime information
- Server statistics
- Environment variables
- Limited system information

## Security Considerations

- **Read-Only Access**: Commands only read system information
- **No Modification**: No commands modify system state
- **User Restrictions**: Can be configured to allow only specific users
- **Command Logging**: All commands are logged for audit purposes

## Customization

You can modify the allowed users in `main.go`:
```go
allowedUsers: map[string]bool{
    "admin": true,
    "sysadmin": true,
    "monitor": true,
},
```

Or add custom administrative commands by extending the `Execute` method.

## Use Cases

- Remote system monitoring
- Server health checks
- Troubleshooting assistance
- System administration over SSH
- Automated monitoring scripts
- DevOps dashboard integration
