#!/bin/bash

# Run Askeladden Beta with beta configuration
echo "🧪 Starting Askeladden Beta..."
echo "Using config: config-beta.yaml"
echo "Using secrets: secrets-beta.yaml"

# Check if beta config files exist
if [ ! -f "config-beta.yaml" ]; then
    echo "❌ Error: config-beta.yaml not found!"
    exit 1
fi

if [ ! -f "secrets-beta.yaml" ]; then
    echo "❌ Error: secrets-beta.yaml not found!"
    exit 1
fi

# Backup original files if they exist
if [ -f "config.yaml" ]; then
    cp config.yaml config.yaml.backup
    echo "📦 Backed up original config.yaml"
fi

if [ -f "secrets.yaml" ]; then
    cp secrets.yaml secrets.yaml.backup
    echo "📦 Backed up original secrets.yaml"
fi

# Copy beta files to expected names
cp config-beta.yaml config.yaml
cp secrets-beta.yaml secrets.yaml
echo "✅ Copied beta configuration files"

# Show which channels will be used
echo "🏗️  Beta configuration:"
echo "   🔧 Prefix: ? (beta) vs ! (production)"
echo "   📋 Log: 1402262636782944366 (bothagen)"
echo "   💬 Main: 1402262679745462453 (kvardagsprat)"
echo "   ⭐ Starboard: 1402262710279864370 (stjernebrettet)"
echo "   ❓ Approval: 1402262743779774568 (spørsmål)"
echo "   💾 Database: daily_questions_testing (isolated from production)"
echo ""

# Run the beta bot
echo "🚀 Starting Askeladden Beta..."
./askeladden-beta

# Clean up - restore original config files
echo "🧹 Restoring original configuration files..."
if [ -f "config.yaml.backup" ]; then
    mv config.yaml.backup config.yaml
    echo "✅ Restored original config.yaml"
else
    rm -f config.yaml
    echo "🗑️  Removed temporary config.yaml"
fi

if [ -f "secrets.yaml.backup" ]; then
    mv secrets.yaml.backup secrets.yaml
    echo "✅ Restored original secrets.yaml"
else
    rm -f secrets.yaml
    echo "🗑️  Removed temporary secrets.yaml"
fi
