#!/bin/bash
set -e

# --- Configuration ---
GO_VERSION="1.22.5"
GO_URL="https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz"
GO_INSTALL_DIR="/usr/local/go"
GO_BINARY="${GO_INSTALL_DIR}/bin/go"
PROJECT_NAME="pivium"
INSTALL_DIR="/opt/${PROJECT_NAME}"
BIN_DIR="/usr/local/bin"

# --- Check for Root ---
if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit
fi

echo "[*] Bootstrapping ${PROJECT_NAME}..."

# --- Install Dependencies ---
echo "[*] Installing base dependencies..."
apt-get update > /dev/null && apt-get install -y git curl > /dev/null

# --- Install Go (Idempotent) ---
if [ ! -f "${GO_BINARY}" ]; then
    echo "[*] Go not found. Installing Go ${GO_VERSION}..."
    curl -sL "${GO_URL}" -o "/tmp/go.tar.gz"
    tar -C "/usr/local" -xzf "/tmp/go.tar.gz"
    rm "/tmp/go.tar.gz"
else
    echo "[*] Go is already installed."
fi
export PATH=$PATH:${GO_INSTALL_DIR}/bin

# --- Setup Project Directory ---
echo "[*] Setting up project directory at ${INSTALL_DIR}"

if [ -d "${INSTALL_DIR}/.git" ]; then
    echo "[*] Git repository found at ${INSTALL_DIR}. Pulling latest changes..."
    cd "${INSTALL_DIR}"
    git pull
else
    echo "[*] Git repository not found. Cloning..."
    # The script is running from inside the repo. We can get the remote URL.
    if [ -d ".git" ]; then
        GIT_REMOTE_URL=$(git config --get remote.origin.url)
        echo "[*] Cloning from ${GIT_REMOTE_URL} to ${INSTALL_DIR}..."
        git clone "${GIT_REMOTE_URL}" "${INSTALL_DIR}"
        cd "${INSTALL_DIR}"
    else
        echo "[!] Error: Not running from a git repository. Cannot clone."
        exit 1
    fi
fi

# --- Build Binaries ---
echo "[*] Building ${PROJECT_NAME}..."
go build -o "${BIN_DIR}/${PROJECT_NAME}" ./cmd/pivium/main.go

echo "[*] Building ${PROJECT_NAME}-updater..."
go build -o "${BIN_DIR}/${PROJECT_NAME}-updater" ./cmd/pivium-updater/main.go

# --- Install Services ---
echo "[*] Installing systemd services..."
bash "${INSTALL_DIR}/scripts/install.sh"

# --- First Run ---
echo "[*] Performing initial reconciliation..."
"${BIN_DIR}/${PROJECT_NAME}" -mode reconcile

echo "[*] Bootstrap complete."