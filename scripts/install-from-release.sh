#!/bin/bash
set -euo pipefail
INSTALL_DIR="${INSTALL_DIR:-/opt/asterisk-response-classifier}"
TARBALL="${1:?install.sh arc-linux-ARCH.tar.gz}"

sudo mkdir -p "$INSTALL_DIR"
sudo tar -xzf "$TARBALL" -C "$INSTALL_DIR"
sudo chmod +x "$INSTALL_DIR/arc" "$INSTALL_DIR/scripts/run-speech-to-phrase.sh"

echo "1) pip install speech-to-phrase"
echo "2) export HASS_TOKEN=... && sudo -E bash $INSTALL_DIR/scripts/run-speech-to-phrase.sh"
echo "3) sudo cp $INSTALL_DIR/deploy/*.service /etc/systemd/system/"
echo "4) sudo systemctl enable --now speech-to-phrase arc"
