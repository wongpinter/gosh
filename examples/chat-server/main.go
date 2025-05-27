package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"repo.nusatek.id/sugeng/gosh"
)

// ChatUser represents a connected user
type ChatUser struct {
	Username  string
	JoinTime  time.Time
	LastSeen  time.Time
	MessageCh chan string
}

// ChatRoom manages chat functionality
type ChatRoom struct {
	users    map[string]*ChatUser
	messages []ChatMessage
	mutex    sync.RWMutex
	maxUsers int
	maxMsgs  int
}

// ChatMessage represents a chat message
type ChatMessage struct {
	Username  string
	Message   string
	Timestamp time.Time
	Type      string // "message", "join", "leave", "system"
}

// NewChatRoom creates a new chat room
func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		users:    make(map[string]*ChatUser),
		messages: make([]ChatMessage, 0),
		maxUsers: 50,
		maxMsgs:  100,
	}
}

// AddUser adds a user to the chat room
func (cr *ChatRoom) AddUser(username string) *ChatUser {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	user := &ChatUser{
		Username:  username,
		JoinTime:  time.Now(),
		LastSeen:  time.Now(),
		MessageCh: make(chan string, 10),
	}

	cr.users[username] = user
	
	// Add join message
	joinMsg := ChatMessage{
		Username:  "System",
		Message:   fmt.Sprintf("%s joined the chat", username),
		Timestamp: time.Now(),
		Type:      "join",
	}
	cr.addMessage(joinMsg)
	
	return user
}

// RemoveUser removes a user from the chat room
func (cr *ChatRoom) RemoveUser(username string) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	if user, exists := cr.users[username]; exists {
		close(user.MessageCh)
		delete(cr.users, username)
		
		// Add leave message
		leaveMsg := ChatMessage{
			Username:  "System",
			Message:   fmt.Sprintf("%s left the chat", username),
			Timestamp: time.Now(),
			Type:      "leave",
		}
		cr.addMessage(leaveMsg)
	}
}

// BroadcastMessage sends a message to all users
func (cr *ChatRoom) BroadcastMessage(msg ChatMessage) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	cr.addMessage(msg)
	
	formattedMsg := cr.formatMessage(msg)
	for _, user := range cr.users {
		select {
		case user.MessageCh <- formattedMsg:
		default:
			// Channel full, skip this user
		}
	}
}

// addMessage adds a message to the history (must be called with lock held)
func (cr *ChatRoom) addMessage(msg ChatMessage) {
	cr.messages = append(cr.messages, msg)
	
	// Keep only the last maxMsgs messages
	if len(cr.messages) > cr.maxMsgs {
		cr.messages = cr.messages[len(cr.messages)-cr.maxMsgs:]
	}
}

// GetUsers returns a list of current users
func (cr *ChatRoom) GetUsers() []string {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	users := make([]string, 0, len(cr.users))
	for username := range cr.users {
		users = append(users, username)
	}
	sort.Strings(users)
	return users
}

// GetRecentMessages returns recent chat messages
func (cr *ChatRoom) GetRecentMessages(count int) []ChatMessage {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	if count > len(cr.messages) {
		count = len(cr.messages)
	}
	
	start := len(cr.messages) - count
	if start < 0 {
		start = 0
	}
	
	return cr.messages[start:]
}

// formatMessage formats a message for display
func (cr *ChatRoom) formatMessage(msg ChatMessage) string {
	timestamp := msg.Timestamp.Format("15:04:05")
	
	switch msg.Type {
	case "join", "leave", "system":
		return fmt.Sprintf("[%s] * %s", timestamp, msg.Message)
	default:
		return fmt.Sprintf("[%s] <%s> %s", timestamp, msg.Username, msg.Message)
	}
}

// Global chat room instance
var chatRoom = NewChatRoom()

// ChatHandler implements the SSH command handler for chat
type ChatHandler struct {
	username string
	user     *ChatUser
}

// NewChatHandler creates a new chat handler for a user
func NewChatHandler(username string) *ChatHandler {
	return &ChatHandler{
		username: username,
	}
}

// Execute implements the CommandHandler interface
func (h *ChatHandler) Execute(cmd string) (string, uint32) {
	cmd = strings.TrimSpace(cmd)
	
	// Initialize user if not done yet
	if h.user == nil {
		h.user = chatRoom.AddUser(h.username)
	}
	
	h.user.LastSeen = time.Now()
	
	if cmd == "" {
		return "", 0
	}
	
	parts := strings.Fields(cmd)
	command := parts[0]
	
	switch command {
	case "/help":
		return h.getHelp(), 0
	case "/users", "/who":
		return h.listUsers(), 0
	case "/history":
		return h.getHistory(parts[1:])
	case "/me":
		return h.sendAction(parts[1:])
	case "/quit", "/exit":
		return "Goodbye! Disconnecting...", 0
	case "/stats":
		return h.getStats(), 0
	case "/time":
		return fmt.Sprintf("Current time: %s", time.Now().Format("2006-01-02 15:04:05")), 0
	default:
		// Regular chat message
		if strings.HasPrefix(cmd, "/") {
			return fmt.Sprintf("Unknown command: %s\nType /help for available commands", command), 1
		}
		
		// Send chat message
		msg := ChatMessage{
			Username:  h.username,
			Message:   cmd,
			Timestamp: time.Now(),
			Type:      "message",
		}
		
		chatRoom.BroadcastMessage(msg)
		return "", 0
	}
}

func (h *ChatHandler) getHelp() string {
	return `Chat Commands:
/help                Show this help message
/users, /who         List online users
/history [count]     Show recent messages (default: 10)
/me <action>         Send an action message
/stats               Show chat statistics
/time                Show current time
/quit, /exit         Leave the chat

To send a message, just type it and press Enter.
Messages starting with / are treated as commands.`
}

func (h *ChatHandler) listUsers() string {
	users := chatRoom.GetUsers()
	if len(users) == 0 {
		return "No users online"
	}
	
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Online users (%d):\n", len(users)))
	for _, user := range users {
		if user == h.username {
			result.WriteString(fmt.Sprintf("  %s (you)\n", user))
		} else {
			result.WriteString(fmt.Sprintf("  %s\n", user))
		}
	}
	
	return result.String()
}

func (h *ChatHandler) getHistory(args []string) (string, uint32) {
	count := 10
	if len(args) > 0 {
		if n, err := fmt.Sscanf(args[0], "%d", &count); n != 1 || err != nil {
			return "Usage: /history [count]", 1
		}
		if count > 50 {
			count = 50
		}
	}
	
	messages := chatRoom.GetRecentMessages(count)
	if len(messages) == 0 {
		return "No messages in history", 0
	}
	
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Last %d messages:\n", len(messages)))
	for _, msg := range messages {
		result.WriteString(chatRoom.formatMessage(msg) + "\n")
	}
	
	return result.String(), 0
}

func (h *ChatHandler) sendAction(args []string) (string, uint32) {
	if len(args) == 0 {
		return "Usage: /me <action>", 1
	}
	
	action := strings.Join(args, " ")
	msg := ChatMessage{
		Username:  h.username,
		Message:   fmt.Sprintf("* %s %s", h.username, action),
		Timestamp: time.Now(),
		Type:      "action",
	}
	
	chatRoom.BroadcastMessage(msg)
	return "", 0
}

func (h *ChatHandler) getStats() string {
	users := chatRoom.GetUsers()
	messages := chatRoom.GetRecentMessages(1000) // Get more for stats
	
	var result strings.Builder
	result.WriteString("=== CHAT STATISTICS ===\n")
	result.WriteString(fmt.Sprintf("Online users: %d\n", len(users)))
	result.WriteString(fmt.Sprintf("Total messages: %d\n", len(messages)))
	
	if len(messages) > 0 {
		oldest := messages[0].Timestamp
		newest := messages[len(messages)-1].Timestamp
		duration := newest.Sub(oldest)
		result.WriteString(fmt.Sprintf("Chat duration: %v\n", duration.Round(time.Second)))
	}
	
	return result.String()
}

// GetPrompt implements the CommandHandler interface
func (h *ChatHandler) GetPrompt() string {
	return fmt.Sprintf("[%s] ", h.username)
}

// GetWelcomeMessage implements the CommandHandler interface
func (h *ChatHandler) GetWelcomeMessage() string {
	users := chatRoom.GetUsers()
	return fmt.Sprintf("Welcome to the Chat Server, %s!\n"+
		"There are currently %d users online.\n"+
		"Type /help for commands or just start chatting!\n"+
		"Type /users to see who's online.",
		h.username, len(users))
}

func main() {
	// Create configuration
	config := sshserver.DefaultConfig()
	config.ListenAddress = ":2226"
	config.HostKeyFile = "server_key"
	config.AuthorizedKeysFile = "authorized_keys"
	config.LogWriter.FilePath = "chat_server.log"

	// Create a custom server that creates different handlers per connection
	server, err := sshserver.NewServer(config, NewChatHandler("default"))
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Chat Server started on port 2226!")
	log.Println("Connect with: ssh -p 2226 <username>@localhost")
	log.Println("Multiple users can connect simultaneously for chat")

	// Wait for interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down chat server...")
	if err := server.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	}
}
