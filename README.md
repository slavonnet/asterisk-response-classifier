# asterisk-response-classifier

[![CI](https://github.com/slavonnet/asterisk-response-classifier/actions/workflows/ci.yml/badge.svg)](https://github.com/slavonnet/asterisk-response-classifier/actions/workflows/ci.yml)

Один сервис на сервере Asterisk: **positive / negative / uncertain**.

## Как speech-to-phrase (но проще)

[speech-to-phrase](https://github.com/OHF-voice/speech-to-phrase) **не переобучает нейросеть**. Она:

1. Берёт **список известных фраз** (шаблоны + имена из HA)
2. Собирает **грамматику** (Kaldi FST) — «какую из *моих* фраз сказали?»
3. При новых сущностях HA **пересобирает грамматику** (секунды, не training с нуля)

Здесь то же самое для да/нет:

| speech-to-phrase | arc |
|------------------|-----|
| Список фраз из HA | `config/phrases.yaml` |
| Пересборка при изменениях | Грамматика Vosk пересоздаётся **на каждый ответ** из yaml |
| Акустическая модель (скачана один раз) | `model/` **уже в Release** |

Открытый STT («что угодно сказал человек») не нужен — только ваши «да» и «нет».

## Установка с Release (всё в одном архиве)

```bash
# скачать arc-linux-arm64.tar.gz (или amd64) из Releases
sudo bash install.sh arc-linux-arm64.tar.gz
```

Внутри tarball: `arc`, `lib/libvosk.so`, **`model/`**, `config/phrases.yaml`, systemd unit.  
**На сервере ничего качать и переобучать не нужно.**

## Dialplan

```asterisk
same => n,Gosub(yesno-ask,s,1(custom/вопрос))
same => n,GotoIf($["${GOSUB_RETVAL}" = "positive"]?yes)
same => n,GotoIf($["${GOSUB_RETVAL}" = "negative"]?no)
same => n,Goto(operator)
```

## Правка фраз

Редактируете `config/phrases.yaml` — списки `positive` / `negative`.  
Перезапуск не нужен: на следующем ответе подхватится новая грамматика.

## Releases

https://github.com/slavonnet/asterisk-response-classifier/releases
