package sshserver

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	"golang.org/x/crypto/ssh"
)

// CommandHandler defines the interface for handling SSH commands
type CommandHandler interface {
	// Execute handles a command and returns the output and exit status
	Execute(cmd string) (string, uint32)
	// GetPrompt returns the shell prompt string
	GetPrompt() string
	// GetWelcomeMessage returns the message shown when a shell session starts
	GetWelcomeMessage() string
}

// Server represents an SSH server instance
type Server struct {
	config        *Config
	sshConfig     *ssh.ServerConfig
	cmdHandler    CommandHandler
	listener      net.Listener
	done         chan struct{}
	wg           sync.WaitGroup
	logger       *log.Logger
}

// NewServer creates a new SSH server instance
func NewServer(config *Config, handler CommandHandler) (*Server, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %v", err)
	}

	var logWriter io.Writer
	if config.LogWriter.Enabled {
		writers := make([]io.Writer, 0)
		
		if config.LogWriter.LogToStdout {
			writers = append(writers, os.Stdout)
		}
		
		if config.LogWriter.FilePath != "" {
			logFile, err := os.OpenFile(config.LogWriter.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, fmt.Errorf("failed to open log file: %v", err)
			}
			writers = append(writers, logFile)
		}
		
		logWriter = io.MultiWriter(writers...)
	}

	s := &Server{
		config:     config,
		cmdHandler: handler,
		done:      make(chan struct{}),
		logger:    log.New(logWriter, "", log.Ldate|log.Ltime|log.Lshortfile),
	}

	sshConfig := &ssh.ServerConfig{
		NoClientAuth: config.NoClientAuth,
	}

	if !config.NoClientAuth {
		private, err := loadHostKey(config.HostKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load host key: %v", err)
		}
		sshConfig.AddHostKey(private)

		sshConfig.PublicKeyCallback = func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			return s.validatePublicKey(conn, key)
		}

		if config.AllowKeyboardInteractive {
			sshConfig.KeyboardInteractiveCallback = s.handleKeyboardInteractive
		}
	}

	s.sshConfig = sshConfig
	return s, nil
}

// Start begins listening for SSH connections
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.config.ListenAddress)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", s.config.ListenAddress, err)
	}

	s.listener = listener
	s.logger.Printf("SSH server listening on %s", s.config.ListenAddress)

	s.wg.Add(1)
	go s.acceptConnections()

	return nil
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	close(s.done)
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			return fmt.Errorf("error closing listener: %v", err)
		}
	}
	s.wg.Wait()
	return nil
}

func (s *Server) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.done:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				select {
				case <-s.done:
					return
				default:
					s.logger.Printf("Failed to accept connection: %v", err)
					continue
				}
			}

			s.wg.Add(1)
			go func() {
				defer s.wg.Done()
				s.handleConnection(conn)
			}()
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	s.logger.Printf("New connection from %s", conn.RemoteAddr())

	sshConn, chans, reqs, err := ssh.NewServerConn(conn, s.sshConfig)
	if err != nil {
		s.logger.Printf("Failed to handshake: %v", err)
		return
	}
	defer sshConn.Close()

	s.logger.Printf("Connection established from %s (user: %s)", sshConn.RemoteAddr(), sshConn.User())

	go s.handleGlobalRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			s.logger.Printf("Could not accept channel: %v", err)
			continue
		}

		go s.handleChannel(channel, requests)
	}
}

func (s *Server) handleChannel(channel ssh.Channel, requests <-chan *ssh.Request) {
	defer channel.Close()

	for req := range requests {
		s.logger.Printf("Received channel request: %s", req.Type)

		switch req.Type {
		case "pty-req":
			req.Reply(true, nil)
		case "shell":
			req.Reply(true, nil)
			if s.cmdHandler != nil {
				channel.Write([]byte(s.cmdHandler.GetWelcomeMessage() + "\n"))
				go s.handleShell(channel)
			}
		case "exec":
			if s.cmdHandler == nil {
				req.Reply(false, nil)
				continue
			}

			command, err := parseExecPayload(req.Payload)
			if err != nil {
				s.logger.Printf("Error parsing exec payload: %v", err)
				req.Reply(false, nil)
				continue
			}

			output, exitStatus := s.cmdHandler.Execute(command)
			channel.Write([]byte(output + "\n"))
			req.Reply(true, nil)
			sendExitStatus(channel, exitStatus)
			return
		default:
			req.Reply(false, nil)
		}
	}
}

func (s *Server) handleShell(channel ssh.Channel) {
	defer channel.Close()

	buffer := make([]byte, 1024)
	var cmdBuffer []byte

	// Send initial prompt
	channel.Write([]byte(s.cmdHandler.GetPrompt()))

	for {
		n, err := channel.Read(buffer)
		if err != nil {
			if err != io.EOF {
				s.logger.Printf("Error reading from channel: %v", err)
			}
			return
		}

		for i := 0; i < n; i++ {
			switch buffer[i] {
			case '\r', '\n':
				if len(cmdBuffer) > 0 {
					cmd := string(cmdBuffer)
					output, _ := s.cmdHandler.Execute(cmd)
					channel.Write([]byte("\r\n" + output + "\r\n" + s.cmdHandler.GetPrompt()))
					cmdBuffer = cmdBuffer[:0]
				} else {
					channel.Write([]byte("\r\n" + s.cmdHandler.GetPrompt()))
				}
			case 0x7f, 0x08: // Backspace
				if len(cmdBuffer) > 0 {
					cmdBuffer = cmdBuffer[:len(cmdBuffer)-1]
					channel.Write([]byte{0x08, 0x20, 0x08}) // Backspace, space, backspace
				}
			default:
				cmdBuffer = append(cmdBuffer, buffer[i])
				channel.Write([]byte{buffer[i]})
			}
		}
	}
}

func (s *Server) handleGlobalRequests(reqs <-chan *ssh.Request) {
	for req := range reqs {
		s.logger.Printf("Received global request: %v", req.Type)
		if req.WantReply {
			req.Reply(false, nil)
		}
	}
}

func (s *Server) validatePublicKey(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	authorizedKeysBytes, err := os.ReadFile(s.config.AuthorizedKeysFile)
	if err != nil {
		s.logger.Printf("Failed to load authorized_keys: %v", err)
		return nil, err
	}

	keyFingerprint := ssh.FingerprintSHA256(key)
	s.logger.Printf("Attempting to authenticate user %s with key %s", conn.User(), keyFingerprint)

	for len(authorizedKeysBytes) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			s.logger.Printf("Error parsing authorized key: %v", err)
			return nil, err
		}

		if ssh.FingerprintSHA256(pubKey) == keyFingerprint {
			s.logger.Printf("Public key authentication successful for user: %s", conn.User())
			return &ssh.Permissions{
				Extensions: map[string]string{
					"pubkey-fp": keyFingerprint,
				},
			}, nil
		}

		authorizedKeysBytes = rest
	}

	return nil, fmt.Errorf("public key authentication failed for %q", conn.User())
}

func (s *Server) handleKeyboardInteractive(conn ssh.ConnMetadata, client ssh.KeyboardInteractiveChallenge) (*ssh.Permissions, error) {
	s.logger.Printf("Keyboard interactive auth attempt from user %s", conn.User())
	return nil, fmt.Errorf("keyboard-interactive authentication not supported")
}

func loadHostKey(keyFile string) (ssh.Signer, error) {
	privateBytes, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %v", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return private, nil
}

func parseExecPayload(payload []byte) (string, error) {
	if len(payload) < 4 {
		return "", fmt.Errorf("exec payload too short")
	}

	length := uint32(payload[3]) | uint32(payload[2])<<8 | uint32(payload[1])<<16 | uint32(payload[0])<<24
	if length == 0 || int(length) > len(payload)-4 {
		return "", fmt.Errorf("invalid command length")
	}

	return string(payload[4 : 4+length]), nil
}

func sendExitStatus(channel ssh.Channel, status uint32) {
	channel.SendRequest("exit-status", false, ssh.Marshal(struct{ Status uint32 }{status}))
}