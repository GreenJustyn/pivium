#!/bin/bash
set -e

SERVICE_FILE="/etc/systemd/system/noded.service"
TIMER_FILE="/etc/systemd/system/noded.timer"
BIN_PATH="/opt/noded/cmd/noded/noded_bin" # In prod, this would be /usr/local/bin/noded

echo "[*] Installing Systemd Units..."

# 1. Create Service Unit
cat <<EOF > $SERVICE_FILE
[Unit]
Description=Noded Infrastructure Reconciler
Wants=network-online.target
After=network-online.target

[Service]
Type=oneshot
WorkingDirectory=/opt/noded
ExecStart=/usr/local/bin/noded -root /opt/noded -mode reconcile
StandardOutput=journal
StandardError=journal
EOF

# 2. Create Timer Unit (Runs every 5 minutes)
cat <<EOF > $TIMER_FILE
[Unit]
Description=Run Noded every 5 minutes

[Timer]
OnBootSec=5min
OnUnitActiveSec=5min

[Install]
WantedBy=timers.target
EOF

# 3. Reload
systemctl daemon-reload
systemctl enable --now noded.timer

echo "[*] Noded scheduled."