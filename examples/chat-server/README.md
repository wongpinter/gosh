# Chat Server SSH Example

This example demonstrates a multi-user chat system accessible via SSH, allowing multiple users to connect simultaneously and chat in real-time.

## Features

- Multi-user real-time chat
- User join/leave notifications
- Chat history
- User list
- Action messages (/me command)
- Chat statistics
- Concurrent connection support

## Setup

1. Generate SSH keys:
   ```bash
   ./setup.sh
   ```

2. Run the server:
   ```bash
   go run main.go
   ```

3. Connect multiple users (from different terminals):
   ```bash
   ssh -p 2226 alice@localhost
   ssh -p 2226 bob@localhost
   ssh -p 2226 charlie@localhost
   ```

## Chat Commands

### Basic Commands
- Just type a message and press Enter to chat
- `/help` - Show available commands
- `/quit`, `/exit` - Leave the chat

### User Management
- `/users`, `/who` - List online users
- `/me <action>` - Send an action message

### History & Information
- `/history [count]` - Show recent messages (default: 10)
- `/stats` - Show chat statistics
- `/time` - Show current time

## Example Chat Session

**Terminal 1 (Alice):**
```
$ ssh -p 2226 alice@localhost
Welcome to the Chat Server, alice!
There are currently 0 users online.
Type /help for commands or just start chatting!
Type /users to see who's online.

[alice] Hello everyone!
[15:30:15] <alice> Hello everyone!
[15:30:22] * bob joined the chat
[alice] /users
Online users (2):
  alice (you)
  bob
[alice] Hi Bob!
[15:30:35] <alice> Hi Bob!
[15:30:40] <bob> Hey Alice! How's it going?
[alice] /me waves
[15:30:45] * alice waves
```

**Terminal 2 (Bob):**
```
$ ssh -p 2226 bob@localhost
Welcome to the Chat Server, bob!
There are currently 1 users online.
Type /help for commands or just start chatting!
Type /users to see who's online.

[15:30:22] <alice> Hello everyone!
[bob] Hey Alice! How's it going?
[15:30:40] <bob> Hey Alice! How's it going?
[15:30:45] * alice waves
[bob] /history 5
Last 4 messages:
[15:30:15] <alice> Hello everyone!
[15:30:35] <alice> Hi Bob!
[15:30:40] <bob> Hey Alice! How's it going?
[15:30:45] * alice waves
```

## Features in Detail

### Real-time Messaging
- Messages are broadcast to all connected users instantly
- Each message includes timestamp and username
- System messages for user join/leave events

### User Management
- Each connection represents a different user
- Username is taken from SSH connection
- Users can see who else is online
- Join/leave notifications

### Chat History
- Server maintains message history (last 100 messages)
- Users can request recent message history
- Timestamps on all messages

### Action Messages
- `/me` command for action-style messages
- Displayed as "* username action"
- Useful for expressing actions or emotions

### Statistics
- Track number of online users
- Message count and chat duration
- Per-user statistics

## Technical Implementation

### Concurrent Connections
- Each SSH connection runs in its own goroutine
- Thread-safe message broadcasting
- Shared chat room state with mutex protection

### Message Broadcasting
- Messages are sent to all connected users
- Non-blocking message delivery
- Channel-based communication

### User Session Management
- Users are added/removed automatically
- Clean disconnection handling
- Session state tracking

## Customization

You can modify the chat server by:

1. **Changing message limits:**
   ```go
   maxMsgs: 500, // Increase message history
   ```

2. **Adding user limits:**
   ```go
   maxUsers: 20, // Limit concurrent users
   ```

3. **Adding chat rooms:**
   - Implement multiple chat rooms
   - Room-based message routing
   - Room switching commands

4. **Adding moderation:**
   - Admin commands
   - Message filtering
   - User management

## Use Cases

- Team communication
- Customer support chat
- Gaming communities
- Educational discussions
- Remote collaboration
- Social networking
