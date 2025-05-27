#!/bin/bash

echo "Setting up Basic SSH Server Example..."

# Generate server host key
if [ ! -f "server_key" ]; then
    echo "Generating server host key..."
    ssh-keygen -t rsa -b 2048 -f server_key -N "" -C "ssh-server-host-key"
    echo "Server host key generated: server_key"
fi

# Generate client key for testing
if [ ! -f "client_key" ]; then
    echo "Generating client key for testing..."
    ssh-keygen -t rsa -b 2048 -f client_key -N "" -C "ssh-client-test-key"
    echo "Client key generated: client_key"
fi

# Create authorized_keys file
echo "Creating authorized_keys file..."
cp client_key.pub authorized_keys
echo "Client public key added to authorized_keys"

# Set proper permissions
chmod 600 server_key client_key
chmod 644 server_key.pub client_key.pub authorized_keys

echo ""
echo "Setup complete!"
echo ""
echo "To start the server:"
echo "  go run main.go"
echo ""
echo "To connect (from another terminal):"
echo "  ssh -i client_key -p 2222 user@localhost"
echo ""
echo "Or without specifying key (if you add client_key to your SSH agent):"
echo "  ssh-add client_key"
echo "  ssh -p 2222 user@localhost"
