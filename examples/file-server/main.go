package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	sshserver "repo.nusatek.id/sugeng/gosh"
)

// FileServerHandler implements a file server over SSH
type FileServerHandler struct {
	rootDir     string
	currentDir  string
	maxFileSize int64 // Maximum file size to display (in bytes)
}

// NewFileServerHandler creates a new file server handler
func NewFileServerHandler(rootDir string) *FileServerHandler {
	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		absRoot = rootDir
	}

	return &FileServerHandler{
		rootDir:     absRoot,
		currentDir:  absRoot,
		maxFileSize: 1024 * 1024, // 1MB
	}
}

// Execute implements the CommandHandler interface
func (h *FileServerHandler) Execute(cmd string) (string, uint32) {
	parts := strings.Fields(strings.TrimSpace(cmd))
	if len(parts) == 0 {
		return "", 0
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "ls", "dir":
		return h.listDirectory(args)
	case "cd":
		return h.changeDirectory(args)
	case "pwd":
		return h.getCurrentDirectory(), 0
	case "cat", "type":
		return h.displayFile(args)
	case "head":
		return h.displayFileHead(args)
	case "tail":
		return h.displayFileTail(args)
	case "stat", "info":
		return h.getFileInfo(args)
	case "find":
		return h.findFiles(args)
	case "download":
		return h.downloadFile(args)
	case "tree":
		return h.showTree(args)
	case "help":
		return h.getHelp(), 0
	default:
		return fmt.Sprintf("Unknown command: %s\nType 'help' for available commands", command), 1
	}
}

func (h *FileServerHandler) listDirectory(args []string) (string, uint32) {
	targetDir := h.currentDir
	if len(args) > 0 {
		targetDir = h.resolvePath(args[0])
	}

	// Security check
	if !h.isPathAllowed(targetDir) {
		return "Error: Access denied", 1
	}

	entries, err := os.ReadDir(targetDir)
	if err != nil {
		return fmt.Sprintf("Error reading directory: %v", err), 1
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Directory: %s\n\n", h.getRelativePath(targetDir)))

	// Sort entries: directories first, then files
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		var typeChar string
		var size string

		if entry.IsDir() {
			typeChar = "d"
			size = "<DIR>"
		} else {
			typeChar = "-"
			size = h.formatFileSize(info.Size())
		}

		modTime := info.ModTime().Format("2006-01-02 15:04")
		result.WriteString(fmt.Sprintf("%s %10s %s %s\n",
			typeChar, size, modTime, entry.Name()))
	}

	return result.String(), 0
}

func (h *FileServerHandler) changeDirectory(args []string) (string, uint32) {
	if len(args) == 0 {
		h.currentDir = h.rootDir
		return fmt.Sprintf("Changed to root directory: %s", h.getRelativePath(h.currentDir)), 0
	}

	targetDir := h.resolvePath(args[0])

	if !h.isPathAllowed(targetDir) {
		return "Error: Access denied", 1
	}

	info, err := os.Stat(targetDir)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), 1
	}

	if !info.IsDir() {
		return "Error: Not a directory", 1
	}

	h.currentDir = targetDir
	return fmt.Sprintf("Changed directory to: %s", h.getRelativePath(h.currentDir)), 0
}

func (h *FileServerHandler) getCurrentDirectory() string {
	return fmt.Sprintf("Current directory: %s", h.getRelativePath(h.currentDir))
}

func (h *FileServerHandler) displayFile(args []string) (string, uint32) {
	if len(args) == 0 {
		return "Usage: cat <filename>", 1
	}

	filePath := h.resolvePath(args[0])
	if !h.isPathAllowed(filePath) {
		return "Error: Access denied", 1
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), 1
	}

	if info.IsDir() {
		return "Error: Cannot display directory content", 1
	}

	if info.Size() > h.maxFileSize {
		return fmt.Sprintf("Error: File too large (%s). Use 'head' or 'tail' instead.",
			h.formatFileSize(info.Size())), 1
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err), 1
	}

	return string(content), 0
}

func (h *FileServerHandler) displayFileHead(args []string) (string, uint32) {
	if len(args) == 0 {
		return "Usage: head <filename> [lines]", 1
	}

	lines := 10
	if len(args) > 1 {
		if n, err := fmt.Sscanf(args[1], "%d", &lines); n != 1 || err != nil {
			return "Error: Invalid line count", 1
		}
	}

	return h.displayFileLines(args[0], lines, true)
}

func (h *FileServerHandler) displayFileTail(args []string) (string, uint32) {
	if len(args) == 0 {
		return "Usage: tail <filename> [lines]", 1
	}

	lines := 10
	if len(args) > 1 {
		if n, err := fmt.Sscanf(args[1], "%d", &lines); n != 1 || err != nil {
			return "Error: Invalid line count", 1
		}
	}

	return h.displayFileLines(args[0], lines, false)
}

func (h *FileServerHandler) displayFileLines(filename string, lineCount int, fromStart bool) (string, uint32) {
	filePath := h.resolvePath(filename)
	if !h.isPathAllowed(filePath) {
		return "Error: Access denied", 1
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), 1
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err), 1
	}

	lines := strings.Split(string(content), "\n")

	var result []string
	if fromStart {
		// Head: take first N lines
		end := lineCount
		if end > len(lines) {
			end = len(lines)
		}
		result = lines[:end]
	} else {
		// Tail: take last N lines
		start := len(lines) - lineCount
		if start < 0 {
			start = 0
		}
		result = lines[start:]
	}

	return strings.Join(result, "\n"), 0
}

func (h *FileServerHandler) getFileInfo(args []string) (string, uint32) {
	if len(args) == 0 {
		return "Usage: stat <filename>", 1
	}

	filePath := h.resolvePath(args[0])
	if !h.isPathAllowed(filePath) {
		return "Error: Access denied", 1
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), 1
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("File: %s\n", h.getRelativePath(filePath)))
	result.WriteString(fmt.Sprintf("Size: %s (%d bytes)\n", h.formatFileSize(info.Size()), info.Size()))
	result.WriteString(fmt.Sprintf("Type: %s\n", h.getFileType(info)))
	result.WriteString(fmt.Sprintf("Modified: %s\n", info.ModTime().Format("2006-01-02 15:04:05")))
	result.WriteString(fmt.Sprintf("Permissions: %s\n", info.Mode().String()))

	return result.String(), 0
}

func (h *FileServerHandler) downloadFile(args []string) (string, uint32) {
	if len(args) == 0 {
		return "Usage: download <filename>", 1
	}

	filePath := h.resolvePath(args[0])
	if !h.isPathAllowed(filePath) {
		return "Error: Access denied", 1
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), 1
	}

	if info.IsDir() {
		return "Error: Cannot download directory", 1
	}

	if info.Size() > h.maxFileSize {
		return fmt.Sprintf("Error: File too large (%s) for download",
			h.formatFileSize(info.Size())), 1
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err), 1
	}

	encoded := base64.StdEncoding.EncodeToString(content)
	return fmt.Sprintf("File: %s\nSize: %s\nBase64 Content:\n%s",
		filepath.Base(filePath), h.formatFileSize(info.Size()), encoded), 0
}

func (h *FileServerHandler) findFiles(args []string) (string, uint32) {
	if len(args) == 0 {
		return "Usage: find <pattern>", 1
	}

	pattern := args[0]
	var matches []string

	err := filepath.Walk(h.currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if !h.isPathAllowed(path) {
			return nil
		}

		name := filepath.Base(path)
		if matched, _ := filepath.Match(pattern, name); matched {
			relPath := h.getRelativePath(path)
			if info.IsDir() {
				matches = append(matches, relPath+"/")
			} else {
				matches = append(matches, relPath)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Sprintf("Error during search: %v", err), 1
	}

	if len(matches) == 0 {
		return fmt.Sprintf("No files found matching pattern: %s", pattern), 0
	}

	result := fmt.Sprintf("Found %d matches for pattern '%s':\n", len(matches), pattern)
	for _, match := range matches {
		result += fmt.Sprintf("  %s\n", match)
	}

	return result, 0
}

func (h *FileServerHandler) showTree(args []string) (string, uint32) {
	targetDir := h.currentDir
	if len(args) > 0 {
		targetDir = h.resolvePath(args[0])
	}

	if !h.isPathAllowed(targetDir) {
		return "Error: Access denied", 1
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Directory tree: %s\n", h.getRelativePath(targetDir)))

	h.buildTree(targetDir, "", &result, 0, 3) // Max depth of 3
	return result.String(), 0
}

func (h *FileServerHandler) buildTree(dir, prefix string, result *strings.Builder, depth, maxDepth int) {
	if depth >= maxDepth {
		return
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for i, entry := range entries {
		isLast := i == len(entries)-1

		var connector string
		if isLast {
			connector = "└── "
		} else {
			connector = "├── "
		}

		result.WriteString(fmt.Sprintf("%s%s%s", prefix, connector, entry.Name()))
		if entry.IsDir() {
			result.WriteString("/")
		}
		result.WriteString("\n")

		if entry.IsDir() && depth < maxDepth-1 {
			var newPrefix string
			if isLast {
				newPrefix = prefix + "    "
			} else {
				newPrefix = prefix + "│   "
			}

			subDir := filepath.Join(dir, entry.Name())
			if h.isPathAllowed(subDir) {
				h.buildTree(subDir, newPrefix, result, depth+1, maxDepth)
			}
		}
	}
}

// Helper methods
func (h *FileServerHandler) resolvePath(path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Clean(filepath.Join(h.currentDir, path))
}

func (h *FileServerHandler) isPathAllowed(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	absRoot, err := filepath.Abs(h.rootDir)
	if err != nil {
		return false
	}

	rel, err := filepath.Rel(absRoot, absPath)
	if err != nil {
		return false
	}

	return !strings.HasPrefix(rel, "..")
}

func (h *FileServerHandler) getRelativePath(path string) string {
	rel, err := filepath.Rel(h.rootDir, path)
	if err != nil {
		return path
	}
	if rel == "." {
		return "/"
	}
	return "/" + filepath.ToSlash(rel)
}

func (h *FileServerHandler) formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func (h *FileServerHandler) getFileType(info os.FileInfo) string {
	if info.IsDir() {
		return "Directory"
	}
	return "Regular file"
}

func (h *FileServerHandler) getHelp() string {
	return `File Server Commands:
- ls [dir]               List directory contents
- cd [dir]               Change directory
- pwd                    Show current directory
- cat <file>             Display file contents
- head <file> [lines]    Show first N lines (default: 10)
- tail <file> [lines]    Show last N lines (default: 10)
- stat <file>            Show file information
- find <pattern>         Find files matching pattern
- download <file>        Download file (base64 encoded)
- tree [dir]             Show directory tree
- help                   Show this help message

Navigation:
- Use relative paths: cd subdir, cat ../file.txt
- Use absolute paths: cd /path/to/dir
- Go to root: cd (no arguments)

Security: Access is restricted to the configured root directory.`
}

// GetPrompt implements the CommandHandler interface
func (h *FileServerHandler) GetPrompt() string {
	relPath := h.getRelativePath(h.currentDir)
	return fmt.Sprintf("files:%s> ", relPath)
}

// GetWelcomeMessage implements the CommandHandler interface
func (h *FileServerHandler) GetWelcomeMessage() string {
	return fmt.Sprintf("Welcome to File Server!\n"+
		"Root directory: %s\n"+
		"Type 'help' to see available commands.\n"+
		"Type 'ls' to list files in current directory.", h.rootDir)
}

func main() {
	// Create a sample directory structure for demonstration
	sampleDir := "./sample_files"
	os.MkdirAll(sampleDir, 0755)
	os.MkdirAll(filepath.Join(sampleDir, "documents"), 0755)
	os.MkdirAll(filepath.Join(sampleDir, "scripts"), 0755)

	// Create sample files
	sampleFiles := map[string]string{
		"readme.txt":          "This is a sample file server.\nYou can browse and view files using SSH commands.",
		"documents/notes.txt": "Meeting Notes\n=============\n\n1. Project status\n2. Next steps\n3. Action items",
		"scripts/hello.sh":    "#!/bin/bash\necho \"Hello from file server!\"\ndate",
		"config.json":         "{\n  \"server\": \"file-server\",\n  \"version\": \"1.0\",\n  \"enabled\": true\n}",
	}

	for filename, content := range sampleFiles {
		fullPath := filepath.Join(sampleDir, filename)
		os.WriteFile(fullPath, []byte(content), 0644)
	}

	// Create configuration
	config := sshserver.DefaultConfig()
	config.ListenAddress = ":2224"
	config.HostKeyFile = "server_key"
	config.AuthorizedKeysFile = "authorized_keys"
	config.LogWriter.FilePath = "file_server.log"

	// Create file server handler
	handler := NewFileServerHandler(sampleDir)

	// Create and start server
	server, err := sshserver.NewServer(config, handler)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Printf("File server started on port 2224!")
	log.Printf("Serving files from: %s", sampleDir)
	log.Println("Connect with: ssh -p 2224 user@localhost")

	// Wait for interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down server...")
	if err := server.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	}
}
