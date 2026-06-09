# AGENTS.md

## Задача

Короткий ulaw от Asterisk → `positive` | `negative` | `uncertain`.  
**Без STT.** Сходство звука с эталонами в `config/references.yaml`.

## Не делать

- STT / Vosk / Docker / текстовые фразы / ASR в этом сервисе

## Запуск

`./arc -port 9099 -config config/references.yaml`

## Эталоны

ulaw 8 kHz в `config/refs/`. Hot-reload yaml на каждый ответ.
