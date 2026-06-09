#!/bin/bash
set -euo pipefail

# Быстрый деплой на Pi / Linux рядом с Asterisk
INSTALL_DIR="${INSTALL_DIR:-/opt/asterisk-response-classifier}"
VOSK_URL="${VOSK_URL:-ws://127.0.0.1:2700}"

ARCH=$(uname -m)
case "$ARCH" in
  aarch64|arm64) BIN=arc-linux-arm64 ;;
  armv7l|armv6l) BIN=arc-linux-armv7 ;;
  x86_64|amd64)  BIN=arc-linux-amd64 ;;
  *) echo "unsupported arch: $ARCH"; exit 1 ;;
esac

echo "==> install to $INSTALL_DIR ($BIN)"
sudo mkdir -p "$INSTALL_DIR"/{bin,config}
sudo cp "dist/$BIN" "$INSTALL_DIR/bin/arc" 2>/dev/null || {
  echo "Скачайте $BIN из GitHub Releases в dist/ или укажите путь"
  exit 1
}
sudo chmod +x "$INSTALL_DIR/bin/arc"
sudo cp config/phrases.yaml "$INSTALL_DIR/config/"

echo "==> vosk-server (docker)"
if command -v docker >/dev/null; then
  docker run -d --name vosk-ru --restart unless-stopped -p 2700:2700 alphacep/kaldi-ru:latest || docker start vosk-ru
fi

echo "==> systemd"
sudo cp deploy/arc.service /etc/systemd/system/arc.service
sudo sed -i "s|ExecStart=.*|ExecStart=$INSTALL_DIR/bin/arc -port 9099 -config $INSTALL_DIR/config/phrases.yaml -vosk-url=$VOSK_URL|" /etc/systemd/system/arc.service
sudo systemctl daemon-reload
sudo systemctl enable --now arc

echo "==> asterisk aeap.conf snippet:"
echo "[response-classifier]"
echo "type=client"
echo "codecs=!all,ulaw"
echo "url=ws://127.0.0.1:9099"
echo "protocol=speech_to_text"
echo ""
echo "Done. Test: dial extension 550"
