#!/bin/bash

# Run Askeladden Beta with beta configuration
echo "🧪 Starting Askeladden Beta..."
echo "Using config: config-beta.yaml"
echo "Using secrets: secrets-beta.yaml" 

# For now, copy files to expected names since the config loader expects specific names
cp config-beta.yaml config.yaml
cp secrets-beta.yaml secrets.yaml

# Run the beta bot
./askeladden-beta

# Clean up - restore original config files
git checkout config.yaml secrets.yaml 2>/dev/null || echo "Original files not found in git"
