#!/bin/bash
set -e

LOG_FILE="/var/log/pivium-update.log"
exec &> >(tee -a "$LOG_FILE")

echo "[*] ($(date))] Starting pivium update..."

cd /opt/pivium

# Fetch latest changes from the git repository
echo "[*] Fetching latest changes from git..."
git fetch origin

# Check if there are any changes
LOCAL=$(git rev-parse @)
REMOTE=$(git rev-parse @{u})

if [ "$LOCAL" == "$REMOTE" ]; then
    echo "[*] Already up-to-date."
    exit 0
fi

echo "[*] Changes detected. Pulling new version..."
git pull

# Re-run bootstrap to build, install and reconcile
echo "[*] Re-running bootstrap script..."
/opt/pivium/scripts/bootstrap.sh

echo "[*] ($(date)) Pivium update complete."
