#!/bin/bash
set -e

REPO="https://github.com/your-org/noded" # Replace with actual
DEST="/opt/noded"

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit
fi

echo "[*] Bootstrapping Debian Node..."

# 1. Base Deps
apt-get update && apt-get install -y git golang curl

# 2. Clone/Pull
if [ -d "$DEST" ]; then
    cd $DEST && git pull
else
    git clone $REPO $DEST
fi

# 3. Build Binary
echo "[*] Building Noded..."
cd $DEST
go build -o /usr/local/bin/noded ./cmd/noded
# Also keep a copy in bin/ for the updater to find in the future
mkdir -p bin
cp /usr/local/bin/noded bin/noded

# 4. Install Service
bash $DEST/scripts/install.sh

# 5. Run once
/usr/local/bin/noded -mode reconcile