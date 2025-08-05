#!/bin/bash

# Deployment script for Askeladden Discord bot
# Usage: ./deploy.sh

set -e

REMOTE_HOST="heim.bitraf.no"
REMOTE_USER="${USER}"
REMOTE_DIR="/home/${REMOTE_USER}/askeladden"
SERVICE_NAME="askeladden"

echo "🚀 Deploying Askeladden bot to ${REMOTE_HOST}..."

# Create remote directory
echo "📁 Creating remote directory..."
ssh "${REMOTE_USER}@${REMOTE_HOST}" "mkdir -p ${REMOTE_DIR}"

# Copy binary and configuration files
echo "📦 Copying files..."
scp askeladden-linux "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/askeladden"
scp config.yaml "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/"
scp secrets.yaml "${REMOTE_USER}@${REMOTE_HOST}:${REMOTE_DIR}/"

# Copy systemd service file
echo "⚙️  Installing systemd service..."
scp askeladden.service "${REMOTE_USER}@${REMOTE_HOST}:/tmp/"

# Set up systemd service and start bot
ssh "${REMOTE_USER}@${REMOTE_HOST}" << 'EOF'
    # Make binary executable
    chmod +x ~/askeladden/askeladden
    
    # Install systemd service
    sudo mv /tmp/askeladden.service /etc/systemd/system/
    sudo systemctl daemon-reload
    sudo systemctl enable askeladden
    
    # Stop service if running, then start
    sudo systemctl stop askeladden || true
    sudo systemctl start askeladden
    
    # Show status
    sudo systemctl status askeladden --no-pager
EOF

echo "✅ Deployment complete!"
echo "🔍 To check logs: ssh ${REMOTE_USER}@${REMOTE_HOST} 'sudo journalctl -u askeladden -f'"
echo "🔄 To restart: ssh ${REMOTE_USER}@${REMOTE_HOST} 'sudo systemctl restart askeladden'"
