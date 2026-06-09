# asterisk-response-classifier

[![CI](https://github.com/slavonnet/asterisk-response-classifier/actions/workflows/ci.yml/badge.svg)](https://github.com/slavonnet/asterisk-response-classifier/actions/workflows/ci.yml)

**Простой** классификатор коротких ответов для Asterisk: `positive` | `negative` | `uncertain`.

## Без STT / ASR / текста

Короткие «да» и «нет» в текст **не переводим** — так точность будет низкой.

Сравниваем **звук ответа** с **эталонными ulaw-записями** (cosine similarity по акустическим признакам).  
ASR прикрутите сами отдельно, если понадобится полный текст.

## Установка

```bash
# из Releases
sudo bash install.sh arc-linux-amd64.tar.gz
```

```bash
./arc -port 9099 -config config/references.yaml
```

## Эталоны

`config/references.yaml`:

```yaml
references:
  positive:
    - refs/positive/da.ulaw
  negative:
    - refs/negative/net.ulaw
```

Запишите ulaw 8 kHz **с вашей линии** — как реально звучат «да» и «нет» у абонентов.  
Несколько вариантов в каждой группе. Файл перечитывается на каждый ответ.

## Asterisk

`aeap.conf`: `url=ws://127.0.0.1:9099`

```asterisk
Gosub(yesno-ask,s,1(custom/вопрос))
GotoIf($["${GOSUB_RETVAL}" = "uncertain"]?operator)
```

## ONNX позже

Флаг `-onnx-model` можно добавить для нейросетевого эмбеддинга — сейчас достаточно эталонов.

## Releases

https://github.com/slavonnet/asterisk-response-classifier/releases
