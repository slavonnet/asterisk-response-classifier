# asterisk-response-classifier

[![CI](https://github.com/slavonnet/asterisk-response-classifier/actions/workflows/ci.yml/badge.svg)](https://github.com/slavonnet/asterisk-response-classifier/actions/workflows/ci.yml)

**Один Go-сервис на том же сервере, где Asterisk.** Без Docker.

Распознаёт ответ абонента → `positive` | `negative` | `uncertain` (→ оператор).

## Зачем Vosk?

Asterisk отдаёт **звук**, а не текст. Нужен STT (speech-to-text), чтобы понять «да» или «нет».  
Vosk — лёгкий офлайн-движок под CPU; встроен **внутрь `arc`** (один процесс, не отдельный сервис).

Цепочка: **звук → Vosk → текст → списки фраз → positive/negative/uncertain**.

## Запуск (3 команды)

На сервере с Asterisk:

```bash
# 1. Модель речи (~50 MB, один раз)
sudo bash scripts/download-model.sh /opt/asterisk-response-classifier/model

# 2. libvosk + arc (из Release или dist/)
sudo bash scripts/install.sh

# 3. Asterisk — aeap.conf (файл asterisk/aeap.conf):
# url=ws://127.0.0.1:9099
```

Или вручную:

```bash
export LD_LIBRARY_PATH=/opt/asterisk-response-classifier/lib
/opt/asterisk-response-classifier/bin/arc \
  -port 9099 \
  -config /opt/asterisk-response-classifier/config/phrases.yaml \
  -model /opt/asterisk-response-classifier/model
```

## Dialplan

```asterisk
same => n,Gosub(yesno-ask,s,1(custom/ваш-вопрос))
same => n,GotoIf($["${GOSUB_RETVAL}" = "positive"]?yes)
same => n,GotoIf($["${GOSUB_RETVAL}" = "negative"]?no)
same => n,Goto(operator)    ; uncertain
```

Пример: `asterisk/extensions.conf.example`.

## Фразы

`config/phrases.yaml` — универсальные «да» / «нет». Перечитывается без рестарта.

## Сборка на сервере (если Release без vosk)

```bash
# libvosk.so рядом или в /usr/local/lib
go build -tags vosk -o arc ./cmd/arc
```

Release-бинарники для linux собираются с `-tags vosk` и `libvosk.so` в комплекте.

## CI / Releases

https://github.com/slavonnet/asterisk-response-classifier/releases

## Лицензия

Apache-2.0
