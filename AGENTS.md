# AGENTS.md

## Проект

Универсальный **да/нет** классификатор для Asterisk: `positive` | `negative` | `uncertain` → оператор.

**Без дерева диалогов** — любой сценарий вызывает `Gosub(yesno-ask,s,1(prompt))`.

## Стек

- Go AEAP WebSocket → Asterisk `SpeechCreate(response-classifier)`
- Vosk websocket STT (`alphacep/kaldi-ru`, порт 2700) — без CGO в `arc`
- Keyword-классификация по `config/phrases.yaml` (hot-reload на каждый ответ)

## Запуск

```bash
docker run -d -p 2700:2700 alphacep/kaldi-ru:latest
./arc -port 9099 -config config/phrases.yaml -vosk-url=ws://127.0.0.1:2700
```

## Dialplan

```
Gosub(yesno-ask,s,1(custom/question))
GotoIf($["${GOSUB_RETVAL}" = "uncertain"]?operator)
```

## Правка фраз (фаза 3)

Редактировать `config/phrases.yaml` — перезапуск **не нужен**.

## CI

GitHub Actions: test + linux amd64/arm64/armv7. Локальный Go не нужен.
