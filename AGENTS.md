# AGENTS.md

## Принцип (как speech-to-phrase)

Не open STT, а **распознавание из списка фраз** из `phrases.yaml`.
Акустическая модель в `model/` (в Release). Грамматика — на лету из yaml.

## Release bundle

`arc-linux-{amd64,arm64}.tar.gz` содержит: arc, libvosk.so, model/, config/, install.sh

## Запуск

```
LD_LIBRARY_PATH=./lib ./arc -model ./model -config ./config/phrases.yaml
```

## Сборка

`go build -tags vosk` — только для разработки; прод = tarball из CI.
