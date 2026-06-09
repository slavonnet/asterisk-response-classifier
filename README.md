# asterisk-response-classifier

[![CI](https://github.com/slavonnet/asterisk-response-classifier/actions/workflows/ci.yml/badge.svg)](https://github.com/slavonnet/asterisk-response-classifier/actions/workflows/ci.yml)

Локальный классификатор ответов абонента для исходящего робота на **Asterisk**.  
Работает на слабом CPU без GPU: **Go + ONNX Runtime + AEAP**.

## Задача

Робот обзванивает с заранее заданными фразами. После каждого вопроса ответ абонента классифицируется:

| Метка | Значение | Действие в dialplan |
|-------|----------|---------------------|
| `positive` | согласие, «да» | следующий узел дерева |
| `negative` | отказ, «нет» | прощание или другая ветка |
| `uncertain` | не распознано / низкая уверенность | **перевод на оператора** |

## Архитектура

```
┌─────────────┐   ulaw audio    ┌──────────────────┐   ONNX (CPU)   ┌─────────────┐
│  Asterisk   │ ──────────────► │  arc (Go AEAP)   │ ─────────────► │ phrase STT  │
│ SpeechCreate│ ◄────────────── │  ws://:9099      │ ◄───────────── │ + classifier│
└─────────────┘  positive/...   └──────────────────┘                └─────────────┘
```

**Почему не C-модуль Asterisk?**  
С Asterisk 18.12+ есть [AEAP](https://docs.asterisk.org/Configuration/Interfaces/Asterisk-External-Application-Protocol-AEAP/) — внешнее приложение по WebSocket. Это тот же `SpeechCreate()`, но логика в Go, без сборки `res_speech_*.so`.

**Почему Go, а не Python/C++?**

- Python — медленнее на слабом CPU, GIL, тяжёлый деплой
- C++ модуль Asterisk — долгая сборка под каждую версию Asterisk
- **Go** — один статический бинарник, AEAP-сервер из коробки, ONNX через готовый `.so`

**Нужны ли wrappers для Go?**  
Да, для ONNX: [`yalue/onnxruntime_go`](https://github.com/yalue/onnxruntime_go) — это CGO-обёртка над **готовой** `libonnxruntime.so`. Сам ONNX Runtime не компилируется.

## CI / релизы

Сборка и тесты — только в **GitHub Actions**, локальный Go не нужен.

| Workflow | Триггер | Что делает |
|----------|---------|------------|
| `ci.yml` | push/PR в `main` | `go test`, `go vet`, сборка linux amd64/arm64/armv7 |
| `release.yml` | тег `v*` | бинарники + GitHub Release |

После первого push CI скачает зависимости сам. Файл `go.sum` можно закоммитить после первого зелёного прогона (Actions → job → скачать или добавить локально через `go mod tidy`, если Go есть).

Релиз:

```bash
git tag v0.1.0
git push origin v0.1.0
```

## Быстрый старт

### 1. Установка бинарника

Скачайте `arc-linux-*` из [GitHub Releases](https://github.com/slavonnet/asterisk-response-classifier/releases) (или артефактов CI) и положите на сервер с Asterisk:

```bash
install -m 755 arc-linux-arm64 /opt/asterisk-response-classifier/bin/arc
```

### 2. Запуск сервиса

```bash
./bin/arc -port 9099 -config config/phrases.yaml
```

### 3. Asterisk

Скопировать `asterisk/aeap.conf` в `/etc/asterisk/aeap.conf` (или include).

Убедиться, что загружены модули:

```
module load res_aeap.so
module load res_speech_aeap.so
```

Фрагмент dialplan — `asterisk/extensions.conf.example`:

```
SpeechCreate(response-classifier)
SpeechStart()
SpeechBackground(custom/greeting,5)
Set(RESPONSE=${SPEECH_TEXT(0)})
```

`RESPONSE` = `positive` | `negative` | `uncertain`.

### 4. Тест

Позвонить на extension `550` (см. example dialplan).

## ONNX / speech-to-phrase

Идея как у [speech-to-phrase](https://github.com/OHF-voice/speech-to-phrase): не «что сказал человек вообще», а **какая из известных фраз** ближе всего. Это на порядки быстрее полного ASR на Raspberry Pi.

План интеграции:

1. **Фаза 1 (сейчас)** — AEAP + keyword-классификатор + dialplan-дерево
2. **Фase 2** — ONNX phrase-matcher / Silero / кастомная модель под ваш словарь
3. **Фаза 3** — переобучение модели при смене скрипта обзвона

Модели класть в `models/` (в `.gitignore`).

## Конфигурация

`config/phrases.yaml`:

- `phrases.positive` / `phrases.negative` — списки фраз
- `tree.nodes` — дерево диалога (документация; маршрутизация в dialplan)

## Требования

- Asterisk 18.12+ / 19.4+ / 21 с `res_aeap`, `res_speech_aeap`
- Linux ARM/x86 для продакшена (бинарник из CI)
- ONNX Runtime shared library для целевой платформы (фаза 2)
## Лицензия

Apache-2.0
