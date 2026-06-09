# AGENTS.md

arc = AEAP → Wyoming → **speech-to-phrase**. Не заменять STP на Vosk/ASR.

- STP: `scripts/run-speech-to-phrase.sh` (tcp://10300)
- arc: `./arc -config config/ivr.yaml`
- Фразы STP: `config/sentences/ru/ivr.yaml`
- Маппинг: `config/ivr.yaml` positive/negative
