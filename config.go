package sshserver

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the SSH server configuration
type Config struct {
	// ListenAddress is the address and port the server listens on (e.g. ":2222")
	ListenAddress string

	// HostKeyFile is the path to the private key used by the server
	HostKeyFile string

	// AuthorizedKeysFile is the path to the authorized_keys file
	AuthorizedKeysFile string

	// NoClientAuth disables client authentication if set to true
	NoClientAuth bool

	// AllowKeyboardInteractive enables keyboard-interactive authentication
	AllowKeyboardInteractive bool

	// LogWriter is where log messages will be written
	LogWriter *LogConfig
}

// LogConfig specifies logging configuration
type LogConfig struct {
	// Enabled turns logging on/off
	Enabled bool

	// FilePath is the path to the log file
	FilePath string

	// LogToStdout determines if logs should also go to stdout
	LogToStdout bool
}

// DefaultConfig returns a new Config with default values
func DefaultConfig() *Config {
	return &Config{
		ListenAddress:      ":2222",
		HostKeyFile:        "server_key",
		AuthorizedKeysFile: "authorized_keys",
		NoClientAuth:       false,
		LogWriter: &LogConfig{
			Enabled:     true,
			FilePath:    "ssh_server.log",
			LogToStdout: true,
		},
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.ListenAddress == "" {
		return fmt.Errorf("listen address cannot be empty")
	}

	if !c.NoClientAuth {
		if c.HostKeyFile == "" {
			return fmt.Errorf("host key file path cannot be empty when client auth is enabled")
		}

		if c.AuthorizedKeysFile == "" {
			return fmt.Errorf("authorized keys file path cannot be empty when client auth is enabled")
		}

		// Check if host key file exists
		if _, err := os.Stat(c.HostKeyFile); err != nil {
			return fmt.Errorf("host key file not found at %s: %v", c.HostKeyFile, err)
		}

		// Check if authorized_keys file exists
		if _, err := os.Stat(c.AuthorizedKeysFile); err != nil {
			return fmt.Errorf("authorized keys file not found at %s: %v", c.AuthorizedKeysFile, err)
		}
	}

	return nil
}

// ResolvePath resolves a relative path to absolute
func ResolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	// Try current working directory first
	cwd, err := os.Getwd()
	if err == nil {
		absPath := filepath.Join(cwd, path)
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	// Try relative to executable
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		absPath := filepath.Join(exeDir, path)
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}

	return path
}
