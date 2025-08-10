#!/bin/bash
set -e

# --- CONFIGURATION ---
SERVER="ellinorlinnea@heim.bitraf.no"
REMOTE_DIR="/home/ellinorlinnea/askeladden"
TMUX_SESSION="askeladden"
REMOTE_LOG_FILE="$REMOTE_DIR/askeladden.log"
REMOTE_APP_PATH="$REMOTE_DIR/askeladden"

LOCAL_BINARY="askeladden-heim"
CONFIG_FILE="config.yaml"
SECRETS_FILE="secrets.yaml"

echo "🚀 Deploying Askeladden..."

# --- LOCAL CHECKS ---
if [ ! -f "$LOCAL_BINARY" ]; then
    echo "❌ Error: $LOCAL_BINARY not found! Please build it first."
    exit 1
fi

# --- REMOTE PREPARATION ---
echo "🔎 Checking remote directory..."
ssh $SERVER "mkdir -p '$REMOTE_DIR'"

# --- FILE TRANSFER ---
echo "📦 Copying files to server..."
scp "$LOCAL_BINARY" "$SERVER:$REMOTE_APP_PATH"
scp "$CONFIG_FILE" "$SECRETS_FILE" "$SERVER:$REMOTE_DIR/"

# --- STOP OLD PROCESS ---
echo "🛑 Stopping old bot in tmux (if any)..."
ssh $SERVER "tmux kill-session -t $TMUX_SESSION 2>/dev/null || true"

# --- SET PERMISSIONS ---
echo "🔑 Making new binary executable..."
ssh $SERVER "chmod +x '$REMOTE_APP_PATH'"

# --- START NEW PROCESS ---
echo "▶️ Starting new tmux session..."
ssh $SERVER "cd '$REMOTE_DIR' && tmux new-session -d -s '$TMUX_SESSION' \"CONFIG_FILE=$CONFIG_FILE SECRETS_FILE=$SECRETS_FILE ./askeladden > '$REMOTE_LOG_FILE' 2>&1\""

# --- VERIFICATION ---
echo "🔎 Verifying deployment..."
sleep 5
if ssh $SERVER "tmux has-session -t '$TMUX_SESSION' 2>/dev/null"; then
    echo "✅ Bot is running in tmux session '$TMUX_SESSION'."
    echo "📋 Recent logs:"
    ssh $SERVER "tail -n 10 '$REMOTE_LOG_FILE'"
else
    echo "❌ CRITICAL: Bot failed to start."
    echo "📋 Full error log:"
    ssh $SERVER "cat '$REMOTE_LOG_FILE' || echo 'No log file found.'"
    exit 1
fi

echo ""
echo "✅ Deployment complete!"
