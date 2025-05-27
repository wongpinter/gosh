# Custom Handler SSH Server Example

This example demonstrates how to create a custom command handler with specialized functionality including calculator, text processing, and utility commands.

## Features

- Custom command handler implementation
- Calculator with basic arithmetic operations
- Random number generation
- Text processing commands (reverse, upper, lower)
- Server statistics tracking
- Time formatting utilities
- Enhanced help system

## Setup

1. Generate SSH keys:
   ```bash
   ./setup.sh
   ```

2. Run the server:
   ```bash
   go run main.go
   ```

3. Connect from another terminal:
   ```bash
   ssh -p 2223 user@localhost
   ```

## Available Commands

### Basic Commands
- `echo <message>` - Echo back the message
- `help` - Show all available commands

### Calculator
- `calc <number1> <operator> <number2>` - Perform calculations
  - Operators: `+`, `-`, `*`, `/`
  - Example: `calc 10 + 5`

### Random Numbers
- `random` - Generate random number (1-100)
- `random <max>` - Generate random number (1-max)
- `random <min> <max>` - Generate random number in range

### Text Processing
- `reverse <text>` - Reverse the text
- `upper <text>` - Convert to uppercase
- `lower <text>` - Convert to lowercase

### Utilities
- `stats` - Show server statistics (uptime, commands executed)
- `time` - Show current time
- `time unix` - Show Unix timestamp
- `time iso` - Show ISO format time
- `time rfc` - Show RFC format time

## Example Session

```
$ ssh -p 2223 user@localhost
Welcome to Custom SSH Server!
This server has enhanced commands for calculations, text processing, and more.
Type 'help' to see all available commands.
custom> calc 15 * 3
15.00 * 3.00 = 45.00
custom> random 1 10
Random number (1-10): 7
custom> reverse Hello World
dlroW olleH
custom> upper hello world
HELLO WORLD
custom> stats
Server Statistics:
- Uptime: 2m15s
- Commands executed: 5
- Started at: 2024-01-15 10:30:45
```

## Implementation Details

The `CustomHandler` struct implements the `CommandHandler` interface with:
- Command parsing and argument handling
- State tracking (command counter, start time)
- Error handling with appropriate exit codes
- Comprehensive help system

This example shows how to extend the SSH server with domain-specific functionality while maintaining clean code organization.
