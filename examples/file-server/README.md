# File Server SSH Example

This example demonstrates an SSH server that provides file browsing and management capabilities, similar to an FTP server but over SSH.

## Features

- Secure file browsing with path restrictions
- Directory navigation (ls, cd, pwd)
- File viewing (cat, head, tail)
- File information and statistics
- File search functionality
- Directory tree visualization
- File download (base64 encoded)
- Automatic sample file generation

## Setup

1. Generate SSH keys and sample files:
   ```bash
   ./setup.sh
   ```

2. Run the server:
   ```bash
   go run main.go
   ```

3. Connect from another terminal:
   ```bash
   ssh -p 2224 user@localhost
   ```

## Available Commands

### Navigation
- `ls [directory]` - List directory contents
- `cd [directory]` - Change directory
- `pwd` - Show current directory
- `tree [directory]` - Show directory tree structure

### File Operations
- `cat <filename>` - Display entire file content
- `head <filename> [lines]` - Show first N lines (default: 10)
- `tail <filename> [lines]` - Show last N lines (default: 10)
- `stat <filename>` - Show detailed file information

### Search and Download
- `find <pattern>` - Find files matching pattern (supports wildcards)
- `download <filename>` - Download file as base64 encoded content

### Utility
- `help` - Show all available commands

## Example Session

```
$ ssh -p 2224 user@localhost
Welcome to File Server!
Root directory: /path/to/sample_files
Type 'help' to see available commands.
Type 'ls' to list files in current directory.

files:/> ls
Directory: /

d      <DIR> 2024-01-15 10:30 documents
d      <DIR> 2024-01-15 10:30 scripts
-      85 B 2024-01-15 10:30 config.json
-      74 B 2024-01-15 10:30 readme.txt

files:/> cd documents
Changed directory to: /documents

files:/documents> ls
Directory: /documents

-      65 B 2024-01-15 10:30 notes.txt

files:/documents> cat notes.txt
Meeting Notes
=============

1. Project status
2. Next steps
3. Action items

files:/documents> cd ..
Changed directory to: /

files:/> find *.txt
Found 2 matches for pattern '*.txt':
  /readme.txt
  /documents/notes.txt

files:/> tree
Directory tree: /
├── config.json
├── documents/
│   └── notes.txt
├── readme.txt
└── scripts/
    └── hello.sh
```

## Security Features

- **Path Restriction**: Access is limited to the configured root directory
- **No Directory Traversal**: Prevents access to parent directories outside the root
- **File Size Limits**: Large files are protected from full display
- **Read-Only Access**: No file modification or deletion capabilities

## Sample Files

The server automatically creates a sample directory structure:
- `readme.txt` - Welcome message
- `config.json` - Sample JSON configuration
- `documents/notes.txt` - Sample meeting notes
- `scripts/hello.sh` - Sample shell script

## Customization

You can modify the root directory in `main.go`:
```go
handler := NewFileServerHandler("/path/to/your/files")
```

Or adjust the maximum file size for display:
```go
handler.maxFileSize = 2 * 1024 * 1024 // 2MB
```
