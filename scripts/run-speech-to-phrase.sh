#!/bin/bash
# Запуск speech-to-phrase — ТОТ ЖЕ движок что на Pi дома.
# Скопируйте пути/токен из вашего рабочего setup Home Assistant.
set -euo pipefail

: "${HASS_TOKEN:?export HASS_TOKEN=...}"
: "${HASS_WS:=ws://127.0.0.1:8123/api/websocket}"

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
MODELS="${MODELS_DIR:-$ROOT/data/models}"
TRAIN="${TRAIN_DIR:-$ROOT/data/train}"
TOOLS="${TOOLS_DIR:-$ROOT/data/tools}"

mkdir -p "$MODELS" "$TRAIN" "$TOOLS"

exec python3 -m speech_to_phrase \
  --uri "tcp://0.0.0.0:10300" \
  --models-dir "$MODELS" \
  --train-dir "$TRAIN" \
  --tools-dir "$TOOLS" \
  --custom-sentences-dir "$ROOT/config/sentences" \
  --hass-token "$HASS_TOKEN" \
  --hass-websocket-uri "$HASS_WS" \
  --retrain-on-start
