#!/bin/bash

echo "Setting up Monitoring Server SSH Example..."

# Generate server host key
if [ ! -f "server_key" ]; then
    echo "Generating server host key..."
    ssh-keygen -t rsa -b 2048 -f server_key -N "" -C "monitoring-server-host-key"
fi

# Generate client key for testing
if [ ! -f "client_key" ]; then
    echo "Generating client key for testing..."
    ssh-keygen -t rsa -b 2048 -f client_key -N "" -C "monitoring-server-client-key"
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
echo "To start the monitoring server:"
echo "  go run main.go"
echo ""
echo "To connect and monitor:"
echo "  ssh -i client_key -p 2228 monitor@localhost"
echo ""
echo "Monitoring commands to try:"
echo "  dashboard    - System overview"
echo "  memory       - Memory metrics"
echo "  runtime      - Runtime info"
echo "  health       - Health check"
echo "  alert        - Check alerts"
echo "  metrics      - List all metrics"
echo "  export json  - Export data"
echo "  help         - Show all commands"
echo ""
echo "The server automatically collects metrics every 30 seconds."
