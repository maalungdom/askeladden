#!/opt/homebrew/bin/fish

# Check if config and secrets files exist
if not test -f "config-beta.yaml"
    echo "Config file (config-beta.yaml) not found."
    exit 1
end

if not test -f "secrets-beta.yaml"
    echo "Secrets file (secrets-beta.yaml) not found."
    exit 1
end

# Run the bot with config and secrets files
set -x CONFIG_FILE config-beta.yaml
set -x SECRETS_FILE secrets-beta.yaml
./askeladden-beta
