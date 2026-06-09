# AGENTS.md

## Проект

Классификатор ответов абонента для исходящего робота Asterisk: **positive / negative / uncertain**.

## Стек

- **Go 1.22** — AEAP WebSocket-сервер (без C-модуля Asterisk)
- **ONNX Runtime** — inference на CPU (Raspberry Pi / слабый сервер)
- **Asterisk 18.12+ / 19.4+ / 21** — `SpeechCreate()` через `res_speech_aeap`

## Структура

```
cmd/arc/              — точка входа
internal/aeap/        — AEAP speech_to_text протокол
internal/classifier/  — positive/negative/uncertain
internal/onnx/        — заглушка под ONNX STT
internal/config/      — phrases.yaml
config/               — фразы и дерево решений
asterisk/             — aeap.conf, extensions.conf.example
```

## Сборка

```bash
go build -o bin/arc ./cmd/arc
./bin/arc -port 9099 -config config/phrases.yaml
```

На Linux ARM (Pi): положить `libonnxruntime.so` рядом или в `/usr/lib`.

## ONNX в Go

Нужен только **wrapper** `github.com/yalue/onnxruntime_go` + **готовый** `libonnxruntime.so`.
Компиляция C++ не требуется — линкуется prebuilt библиотека через CGO.

## Следующие шаги

1. Подключить ONNX STT (phrase-limited, по аналогии с speech-to-phrase)
2. Добавить ulaw→PCM декодер перед inference
3. Обучить/экспортировать маленькую модель под фиксированный набор фраз
4. Systemd unit для `arc` на сервере с Asterisk

## Dialplan

`SPEECH_TEXT(0)` → `positive` | `negative` | `uncertain`
При `uncertain` — перевод на оператора (см. `asterisk/extensions.conf.example`).
