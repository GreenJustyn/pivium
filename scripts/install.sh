#!/bin/bash
set -e

PROJECT_NAME="pivium"
INSTALL_DIR="/opt/${PROJECT_NAME}"

# Service files from the manifests directory
PIVIUM_SERVICE_SRC="${INSTALL_DIR}/manifests/systemd/${PROJECT_NAME}.service"
PIVIUM_TIMER_SRC="${INSTALL_DIR}/manifests/systemd/${PROJECT_NAME}.timer"
UPDATER_SERVICE_SRC="${INSTALL_DIR}/manifests/systemd/${PROJECT_NAME}-updater.service"

# Destination for systemd files
PIVIUM_SERVICE_DST="/etc/systemd/system/${PROJECT_NAME}.service"
PIVIUM_TIMER_DST="/etc/systemd/system/${PROJECT_NAME}.timer"
UPDATER_SERVICE_DST="/etc/systemd/system/${PROJECT_NAME}-updater.service"

echo "[*] Installing Systemd Units for ${PROJECT_NAME}..."

# 1. Install Pivium reconciler service and timer
if [ -f "${PIVIUM_SERVICE_SRC}" ] && [ -f "${PIVIUM_TIMER_SRC}" ]; then
    echo "[*] Installing ${PROJECT_NAME} service and timer..."
    cp "${PIVIUM_SERVICE_SRC}" "${PIVIUM_SERVICE_DST}"
    cp "${PIVIUM_TIMER_SRC}" "${PIVIUM_TIMER_DST}"
    systemctl enable --now "${PROJECT_NAME}.timer"
else
    echo "[!] Warning: ${PROJECT_NAME} service or timer manifest not found. Skipping."
fi

# 2. Install Pivium updater service
if [ -f "${UPDATER_SERVICE_SRC}" ]; then
    echo "[*] Installing ${PROJECT_NAME}-updater service..."
    cp "${UPDATER_SERVICE_SRC}" "${UPDATER_SERVICE_DST}"

    # Create the updater config directory and a placeholder config file
    mkdir -p "/etc/pivium"
    if [ ! -f "/etc/pivium/updater.conf" ]; then
        echo "# Please set your webhook secret here" > /etc/pivium/updater.conf
        echo "PIVIUM_UPDATER_SECRET=your_super_secret_token" >> /etc/pivium/updater.conf
    fi

    systemctl enable --now "${PROJECT_NAME}-updater.service"
else
    echo "[!] Warning: ${PROJECT_NAME}-updater service manifest not found. Skipping."
fi

# 3. Reload systemd
echo "[*] Reloading systemd daemon..."
systemctl daemon-reload

echo "[*] Service installation complete."