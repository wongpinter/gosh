#!/bin/bash

echo "Setting up Admin Panel SSH Server Example..."

# Generate server host key
if [ ! -f "server_key" ]; then
    echo "Generating server host key..."
    ssh-keygen -t rsa -b 2048 -f server_key -N "" -C "admin-panel-host-key"
fi

# Generate client key for testing
if [ ! -f "client_key" ]; then
    echo "Generating client key for testing..."
    ssh-keygen -t rsa -b 2048 -f client_key -N "" -C "admin-panel-client-key"
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
echo "To start the admin server:"
echo "  go run main.go"
echo ""
echo "To connect as admin:"
echo "  ssh -i client_key -p 2225 admin@localhost"
echo ""
echo "Try these administrative commands:"
echo "  status"
echo "  memory"
echo "  processes"
echo "  disk"
echo "  logs 10"
echo "  help"
echo ""
echo "Note: Some commands require Linux and appropriate permissions."
