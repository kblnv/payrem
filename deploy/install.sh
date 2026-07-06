#!/bin/bash
set -e

INSTALL_DIR="/opt/payr"

command -v git >/dev/null 2>&1 || { echo "Git required"; exit 1; }
command -v go >/dev/null 2>&1 || { echo "Go required"; exit 1; }
command -v make >/dev/null 2>&1 || { echo "Make required"; exit 1; }

mkdir -p "$INSTALL_DIR"/{bin,plugins,config}
PAYR_CORE_OUTPUT="$INSTALL_DIR/bin/payr" make core
PAYR_PLUGINS_OUTPUT="$INSTALL_DIR/plugins" make plugins

$INSTALL_DIR/bin/payr init --path $INSTALL_DIR/config/payr.config.json

SYSTEMD_SERVICE='[Unit]
Description=Payr Notification Service

[Service]
Type=simple
ExecStart=/opt/payr/bin/payr run --config /opt/payr/config/payr.config.json
WorkingDirectory=/opt/payr
Restart=on-failure

[Install]
WantedBy=multi-user.target'

echo "$SYSTEMD_SERVICE" > /etc/systemd/system/payr.service
systemctl daemon-reload
systemctl enable payr

echo "Done! Edit $INSTALL_DIR/config/payr.config.json, then: systemctl start payr"
