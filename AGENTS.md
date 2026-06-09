# AGENTS.md

Phrase-limited recognition как speech-to-phrase. Не ASR.

- `config/sentences.yaml` — lists + intents (формат как STP custom sentences)
- `arc` + `-model model/` — acoustic model в Release
- AEAP → `positive|negative|uncertain`

Сборка: `go build -tags phrase`
