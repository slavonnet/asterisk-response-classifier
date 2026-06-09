#!/bin/bash
# Распаковка tar.gz из GitHub Release — всё уже внутри (модель, libvosk, arc).
set -euo pipefail

INSTALL_DIR="${INSTALL_DIR:-/opt/asterisk-response-classifier}"
TARBALL="${1:?usage: install.sh arc-linux-ARCH.tar.gz}"

echo "==> extract to $INSTALL_DIR"
sudo mkdir -p "$INSTALL_DIR"
sudo tar -xzf "$TARBALL" -C "$INSTALL_DIR"

echo "==> systemd"
sudo cp "$INSTALL_DIR/arc.service" /etc/systemd/system/arc.service
sudo sed -i "s|INSTALL_DIR|$INSTALL_DIR|g" /etc/systemd/system/arc.service
sudo systemctl daemon-reload
sudo systemctl enable --now arc

cat <<EOF

Готово. Модель уже в $INSTALL_DIR/model — ничего скачивать не нужно.
Правьте фразы: $INSTALL_DIR/config/phrases.yaml (без рестарта).

Asterisk aeap.conf:
  url=ws://127.0.0.1:9099
EOF
