#!/bin/bash

echo "Setting up Game Server SSH Example..."

# Generate server host key
if [ ! -f "server_key" ]; then
    echo "Generating server host key..."
    ssh-keygen -t rsa -b 2048 -f server_key -N "" -C "game-server-host-key"
fi

# Generate client key for testing
if [ ! -f "client_key" ]; then
    echo "Generating client key for testing..."
    ssh-keygen -t rsa -b 2048 -f client_key -N "" -C "game-server-client-key"
fi

# Create authorized_keys file
echo "Creating authorized_keys file..."
cp client_key.pub authorized_keys

# Set proper permissions
chmod 600 server_key client_key
chmod 644 server_key.pub client_key.pub authorized_keys

echo ""
echo "Setup complete!"
echo ""
echo "To start the game server:"
echo "  go run main.go"
echo ""
echo "To connect and play:"
echo "  ssh -i client_key -p 2227 player@localhost"
echo ""
echo "Available games to try:"
echo "  guess    - Number guessing game"
echo "  rps      - Rock Paper Scissors"
echo "  quiz     - Trivia questions"
echo ""
echo "Game commands:"
echo "  menu     - Main menu"
echo "  score    - Show your score"
echo "  help     - Show help"
echo "  quit     - Exit"
echo ""
echo "Have fun gaming over SSH! 🎮"
