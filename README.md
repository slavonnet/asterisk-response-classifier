# asterisk-response-classifier

[![CI](https://github.com/slavonnet/asterisk-response-classifier/actions/workflows/ci.yml/badge.svg)](https://github.com/slavonnet/asterisk-response-classifier/actions/workflows/ci.yml)

Локальный **да/нет** классификатор для исходящего робота Asterisk.  
CPU-only, без GPU: **Go + Vosk STT + AEAP**.

## Что делает

После любого вопроса робота ответ абонента → одна из меток:

| Метка | Примеры | Действие |
|-------|---------|----------|
| `positive` | да, ага, конечно, хорошо | ваша ветка «да» |
| `negative` | нет, не надо, не интересно | ваша ветка «нет» |
| `uncertain` | непонятно, тишина, постороннее | **оператор** |

**Дерево диалогов не нужно** — один макрос `yesno-ask` на любой вопрос.

## Быстрый деплой (сегодня)

### 1. Vosk STT (Docker)

```bash
docker run -d --name vosk-ru --restart unless-stopped -p 2700:2700 alphacep/kaldi-ru:latest
```

### 2. arc (бинарник из [Releases](https://github.com/slavonnet/asterisk-response-classifier/releases) или CI)

```bash
chmod +x arc-linux-arm64   # Pi
./arc-linux-arm64 -port 9099 -config config/phrases.yaml -vosk-url=ws://127.0.0.1:2700
```

Или всё сразу:

```bash
cd deploy && docker compose up -d
```

### 3. Asterisk

`asterisk/aeap.conf`:

```ini
[response-classifier]
type=client
codecs=!all,ulaw
url=ws://127.0.0.1:9099
protocol=speech_to_text
```

Модули: `res_aeap`, `res_speech_aeap`.

### 4. Dialplan — любой диалог

```asterisk
same => n,Gosub(yesno-ask,s,1(custom/my-question))
same => n,GotoIf($["${GOSUB_RETVAL}" = "positive"]?yes-branch)
same => n,GotoIf($["${GOSUB_RETVAL}" = "negative"]?no-branch)
same => n,Goto(operator)    ; uncertain
```

Полный пример: `asterisk/extensions.conf.example`. Тест: extension `550`.

## Правка фраз без рестарта

`config/phrases.yaml` — списки `positive` / `negative`.  
Файл перечитывается **на каждый ответ**. Меняете YAML → сразу работает.

## Архитектура

```
Asterisk (ulaw) ──AEAP──► arc:9099 ──PCM──► Vosk:2700
                              │
                              └── keyword match → positive|negative|uncertain
```

## CI / релиз

- Push в `main` → тесты + бинарники linux amd64/arm64/armv7
- `git tag v0.2.0 && git push origin v0.2.0` → GitHub Release

## Скрипт установки на Pi

```bash
# положите arc-linux-arm64 в dist/
bash scripts/install.sh
```

## Лицензия

Apache-2.0
