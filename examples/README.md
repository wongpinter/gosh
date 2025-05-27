# SSH Server Examples

This directory contains various examples demonstrating different use cases for the SSH server package.

## Examples Overview

### 1. Basic Server (`basic-server/`)
A simple SSH server with default configuration and basic commands.

### 2. Custom Handler (`custom-handler/`)
Demonstrates how to create custom command handlers with specialized functionality.

### 3. File Server (`file-server/`)
SSH server that allows browsing and downloading files from the server.

### 4. Admin Panel (`admin-panel/`)
Administrative SSH interface with system monitoring and management commands.

### 5. Chat Server (`chat-server/`)
Multi-user chat system accessible via SSH.

### 6. Proxy Server (`proxy-server/`)
SSH server that acts as a proxy or gateway to other services.

### 7. Monitoring Server (`monitoring-server/`)
Server monitoring and metrics collection via SSH interface.

### 8. Game Server (`game-server/`)
Simple text-based games playable over SSH.

## Getting Started

Each example directory contains:
- `main.go` - The main application code
- `README.md` - Specific setup instructions
- `setup.sh` - Script to generate required keys and files
- Configuration files as needed

## Prerequisites

1. Go 1.24.3 or later
2. SSH client for testing
3. OpenSSH tools for key generation

## Quick Start

1. Choose an example directory
2. Run the setup script: `./setup.sh`
3. Build and run: `go run main.go`
4. Connect with SSH client: `ssh -p 2222 user@localhost`

## Security Note

These examples are for demonstration purposes. For production use:
- Use strong host keys
- Implement proper authentication
- Configure appropriate logging
- Follow security best practices
