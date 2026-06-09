# asterisk-response-classifier

**arc** — тонкий мост Asterisk (AEAP) → **speech-to-phrase** (Wyoming).  
Никакого Vosk. Никакого своего ASR. Тот же движок, что у вас на Pi.

## Запуск завтра (2 сервиса на одной машине)

### 1. speech-to-phrase (как дома на Pi)

```bash
pip install speech-to-phrase
export HASS_TOKEN=ваш_токен_HA
sudo bash scripts/run-speech-to-phrase.sh
```

Или systemd: `deploy/speech-to-phrase.service` (пропишите HASS_TOKEN).

Фразы: `config/sentences/ru/ivr.yaml` — формат custom sentences speech-to-phrase.

### 2. arc (бинарник ~2 MB)

```bash
./arc -port 9099 -config config/ivr.yaml
```

Asterisk `aeap.conf`: `url=ws://127.0.0.1:9099`

### 3. Dialplan

```asterisk
Gosub(yesno-ask,s,1(custom/вопрос))
GotoIf($["${GOSUB_RETVAL}" = "uncertain"]?operator)
```

## Схема

```
Asterisk ──AEAP──► arc:9099 ──Wyoming tcp──► speech-to-phrase:10300
                      │                           (Kaldi, как на Pi)
                      └── positive|negative|uncertain
```

## Releases

Только `arc` + конфиги. speech-to-phrase ставите сами (`pip`), как уже работает дома.

https://github.com/slavonnet/asterisk-response-classifier/releases
