# Basic SSH Server Example

This example demonstrates a simple SSH server with default configuration and basic commands.

## Features

- Default SSH server configuration
- Basic command handler with built-in commands
- Public key authentication
- Logging to file and stdout
- Graceful shutdown

## Setup

1. Generate SSH host key and authorized keys:
   ```bash
   ./setup.sh
   ```

2. Build and run the server:
   ```bash
   go run main.go
   ```

3. Connect from another terminal:
   ```bash
   ssh -p 2222 user@localhost
   ```

## Available Commands

Once connected, you can use these commands:
- `hello` - Simple greeting
- `getDate` - Get current server time
- `uptime` - Show server uptime
- `help` - List all available commands

## Files Created

- `server_key` - SSH host private key
- `server_key.pub` - SSH host public key
- `authorized_keys` - Authorized public keys for client authentication
- `client_key` - Client private key (for testing)
- `client_key.pub` - Client public key
- `ssh_server.log` - Server log file

## Testing

After running setup.sh, you can connect using the generated client key:
```bash
ssh -i client_key -p 2222 user@localhost
```

## Customization

You can modify the configuration in `main.go`:
- Change the listen port
- Modify log settings
- Add custom commands to the handler
