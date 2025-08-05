#!/bin/bash

# Deployment script for Askeladden
# Deploys beta branch code as production bot using tmux
# Uses SSH multiplexing to reduce password prompts
set -e

SERVER="ellinorlinnea@heim.bitraf.no"
REMOTE_DIR="/home/ellinorlinnea/askeladden"
LOCAL_BINARY="askeladden-linux"
CONFIG_FILE="config.yaml"
SECRETS_FILE="secrets.yaml"
TMUX_SESSION="askeladden"
SSH_CONTROL="/tmp/askeladden-deploy-%r@%h:%p"

echo "🚀 Deploying Askeladden to production..."

# Check if binary exists
if [ ! -f "$LOCAL_BINARY" ]; then
    echo "❌ Error: $LOCAL_BINARY not found! Please build it first:"
    echo "   GOOS=linux GOARCH=amd64 go build -o askeladden-linux cmd/askeladden/main.go cmd/askeladden/scheduler.go"
    exit 1
fi

# Check if config exists
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ Error: $CONFIG_FILE not found!"
    exit 1
fi

# Check if secrets exists
if [ ! -f "$SECRETS_FILE" ]; then
    echo "❌ Error: $SECRETS_FILE not found!"
    exit 1
fi

# Setup SSH multiplexing
echo "🔗 Establishing SSH connection..."
ssh -M -S "$SSH_CONTROL" -fN $SERVER

# Function to run SSH commands using the multiplexed connection
ssh_run() {
    ssh -S "$SSH_CONTROL" $SERVER "$@"
}

# Function to copy files using the multiplexed connection
scp_copy() {
    scp -o "ControlPath=$SSH_CONTROL" "$@"
}

# Cleanup function to close SSH connection
cleanup() {
    echo "🔌 Closing SSH connection..."
    ssh -S "$SSH_CONTROL" -O exit $SERVER 2>/dev/null || true
}

# Set trap to cleanup on exit
trap cleanup EXIT

echo "🔄 Stopping existing bot processes..."

# Create remote directory if it doesn't exist
ssh_run "mkdir -p $REMOTE_DIR"

# Stop any existing tmux session FIRST
echo "   → Stopping existing tmux session..."
ssh_run "tmux kill-session -t $TMUX_SESSION 2>/dev/null || echo 'No existing session found'"

# Kill any remaining askeladden processes
echo "   → Killing any remaining bot processes..."
ssh_run "pkill -f 'askeladden' || echo 'No processes to kill'"

# Wait for processes to fully stop
echo "   → Waiting for processes to stop..."
sleep 3

echo "📦 Copying files to server..."

# Copy all files in one batch
echo "   → Copying binary..."
scp_copy $LOCAL_BINARY $SERVER:$REMOTE_DIR/askeladden

echo "   → Copying config..."
scp_copy $CONFIG_FILE $SERVER:$REMOTE_DIR/

echo "   → Copying secrets..."
scp_copy $SECRETS_FILE $SERVER:$REMOTE_DIR/

# Make binary executable
ssh_run "chmod +x $REMOTE_DIR/askeladden"

# Wait a moment for processes to stop
sleep 2

# Start new tmux session with the bot
echo "   → Starting new tmux session..."
ssh_run "cd $REMOTE_DIR && tmux new-session -d -s $TMUX_SESSION './askeladden'"

# Wait a moment for the bot to start
sleep 3

# Check if the session is running
echo "   → Checking tmux session..."
if ssh_run "tmux has-session -t $TMUX_SESSION 2>/dev/null"; then
    echo "✅ Bot is running in tmux session '$TMUX_SESSION'"
    
    # Show the first few lines of output
    echo "📋 Recent output:"
    ssh_run "tmux capture-pane -t $TMUX_SESSION -p | tail -10"
else
    echo "❌ Failed to start tmux session"
    exit 1
fi

echo ""
echo "✅ Deployment complete!"
echo ""
echo "🔧 Useful commands:"
echo "   View session: ssh $SERVER 'tmux attach -t $TMUX_SESSION'"
echo "   View logs:    ssh $SERVER 'tmux capture-pane -t $TMUX_SESSION -p'"
echo "   Kill bot:     ssh $SERVER 'tmux kill-session -t $TMUX_SESSION'"
echo "   List tmux:    ssh $SERVER 'tmux list-sessions'"
echo ""
echo "📝 To detach from tmux session when attached: Ctrl+B, then D"
