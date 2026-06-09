#!/bin/bash
# Скачать русскую модель Vosk (~50 MB) — один раз на сервер с Asterisk
set -euo pipefail

MODEL_DIR="${1:-/opt/asterisk-response-classifier/model}"
MODEL_URL="${MODEL_URL:-https://alphacephei.com/vosk/models/vosk-model-small-ru-0.22.zip}"

mkdir -p "$(dirname "$MODEL_DIR")"
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "==> download $MODEL_URL"
curl -fsSL "$MODEL_URL" -o "$tmp/model.zip"
unzip -q "$tmp/model.zip" -d "$tmp"
mv "$tmp"/vosk-model-small-ru-0.22 "$MODEL_DIR"
echo "==> model ready: $MODEL_DIR"
