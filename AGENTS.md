# AGENTS.md

## Суть

Один процесс `arc` на хосте Asterisk. **Без Docker.**

- AEAP WebSocket `:9099` ← Asterisk `SpeechCreate(response-classifier)`
- STT: Vosk in-process (`-model /path`, build `-tags vosk`)
- Классификация: `config/phrases.yaml` (hot-reload)

## Запуск

```bash
LD_LIBRARY_PATH=/opt/asterisk-response-classifier/lib \
  ./arc -port 9099 -config config/phrases.yaml -model /opt/.../model
```

## Dialplan

`Gosub(yesno-ask,s,1(prompt))` → `GOSUB_RETVAL` = positive|negative|uncertain

## Сборка с STT

```bash
go build -tags vosk -o arc ./cmd/arc
```

Без `-tags vosk` — только тесты/CI, STT не работает.
