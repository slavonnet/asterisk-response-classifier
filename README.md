# asterisk-response-classifier

Классификатор коротких ответов для Asterisk — **как [speech-to-phrase](https://github.com/OHF-voice/speech-to-phrase)**:

> не «что человек сказал вообще», а **«какая из известных фраз»**

- `positive` / `negative` / `uncertain` → dialplan
- Конфиг `config/sentences.yaml` — тот же принцип, что custom sentences у speech-to-phrase
- Меняете список фraz → на следующем ответе новая грамматика (без переобучения на сервере)
- **Release tarball**: `arc` + `model/` + `libvosk.so` + `sentences.yaml` — распаковал и работает

## Установка

```bash
sudo bash install.sh arc-linux-arm64.tar.gz
```

## Конфиг (как speech-to-phrase)

```yaml
lists:
  yes_word:
    values:
      - in: "да"
      - in: "ага"
intents:
  Positive:
    data:
      - sentences: ["{yes_word}"]
```

## Asterisk

`aeap.conf`: `url=ws://127.0.0.1:9099`

```asterisk
Gosub(yesno-ask,s,1(custom/вопрос))
GotoIf($["${GOSUB_RETVAL}" = "uncertain"]?operator)
```

ASR для полного текста — отдельно, сюда не входит.

## Releases

https://github.com/slavonnet/asterisk-response-classifier/releases
