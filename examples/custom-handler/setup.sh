#!/bin/bash

echo "Setting up Custom Handler SSH Server Example..."

# Generate server host key
if [ ! -f "server_key" ]; then
    echo "Generating server host key..."
    ssh-keygen -t rsa -b 2048 -f server_key -N "" -C "custom-ssh-server-host-key"
fi

# Generate client key for testing
if [ ! -f "client_key" ]; then
    echo "Generating client key for testing..."
    ssh-keygen -t rsa -b 2048 -f client_key -N "" -C "custom-ssh-client-test-key"
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
echo "To start the server:"
echo "  go run main.go"
echo ""
echo "To connect:"
echo "  ssh -i client_key -p 2223 user@localhost"
echo ""
echo "Try these commands once connected:"
echo "  calc 10 + 5"
echo "  random 1 100"
echo "  reverse Hello World"
echo "  stats"
echo "  help"
