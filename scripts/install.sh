#!/bin/bash
# Установка на тот же Linux-хост, где Asterisk. Без Docker.
set -euo pipefail

INSTALL_DIR="${INSTALL_DIR:-/opt/asterisk-response-classifier}"
MODEL_DIR="${MODEL_DIR:-$INSTALL_DIR/model}"
VOSK_VER="${VOSK_VER:-0.3.45}"

ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64)  VOSK_ARCH=x86_64; BIN=arc-linux-amd64 ;;
  aarch64|arm64) VOSK_ARCH=aarch64; BIN=arc-linux-arm64 ;;
  armv7l|armv6l)  VOSK_ARCH=armhf;  BIN=arc-linux-armv7 ;;
  *) echo "unsupported arch: $ARCH"; exit 1 ;;
esac

echo "==> dirs"
sudo mkdir -p "$INSTALL_DIR"/{bin,config,lib}

echo "==> libvosk"
if ! ldconfig -p 2>/dev/null | grep -q libvosk; then
  tmp=$(mktemp -d)
  curl -fsSL "https://github.com/alphacep/vosk-api/releases/download/v${VOSK_VER}/vosk-linux-${VOSK_ARCH}-${VOSK_VER}.zip" -o "$tmp/vosk.zip"
  unzip -q "$tmp/vosk.zip" -d "$tmp"
  sudo cp "$tmp/vosk-linux-${VOSK_ARCH}-${VOSK_VER}/libvosk.so" "$INSTALL_DIR/lib/"
  sudo ldconfig "$INSTALL_DIR/lib" 2>/dev/null || true
  rm -rf "$tmp"
fi

echo "==> arc binary"
if [ -f "dist/$BIN" ]; then
  sudo cp "dist/$BIN" "$INSTALL_DIR/bin/arc"
elif [ -f "$INSTALL_DIR/bin/arc" ]; then
  echo "using existing $INSTALL_DIR/bin/arc"
else
  echo "Скачайте $BIN из GitHub Releases в dist/ или соберите: go build -tags vosk -o arc ./cmd/arc"
  exit 1
fi
sudo chmod +x "$INSTALL_DIR/bin/arc"
sudo cp config/phrases.yaml "$INSTALL_DIR/config/"

echo "==> speech model"
if [ ! -f "$MODEL_DIR/conf/model.conf" ]; then
  sudo bash scripts/download-model.sh "$MODEL_DIR"
fi

echo "==> systemd"
sudo cp deploy/arc.service /etc/systemd/system/arc.service
sudo sed -i "s|INSTALL_DIR|$INSTALL_DIR|g" /etc/systemd/system/arc.service
sudo sed -i "s|MODEL_DIR|$MODEL_DIR|g" /etc/systemd/system/arc.service
sudo systemctl daemon-reload
sudo systemctl enable --now arc

cat <<EOF

Готово. Один процесс arc на этом же сервере.

Asterisk aeap.conf:
[response-classifier]
type=client
codecs=!all,ulaw
url=ws://127.0.0.1:9099
protocol=speech_to_text

Тест: dial 550
EOF
