#!/bin/bash
set -euo pipefail
INSTALL_DIR="${INSTALL_DIR:-/opt/asterisk-response-classifier}"
TARBALL="${1:?install.sh arc-linux-ARCH.tar.gz}"

sudo mkdir -p "$INSTALL_DIR"
sudo tar -xzf "$TARBALL" -C "$INSTALL_DIR"
sudo cp "$INSTALL_DIR/arc.service" /etc/systemd/system/arc.service
sudo sed -i "s|INSTALL_DIR|$INSTALL_DIR|g" /etc/systemd/system/arc.service
sudo systemctl daemon-reload
sudo systemctl enable --now arc
echo "OK. Правьте фразы: $INSTALL_DIR/config/sentences.yaml"
