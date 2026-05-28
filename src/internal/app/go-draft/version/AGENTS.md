# version — Версионирование приложения

## Назначение

Пакет `version` — единственный источник истины для версии приложения.

## Структура

Единственный файл — `version.go`. Содержит:

- **`Version`** — строка с версией в формате SemVer: `"v0.5.0"`

## Использование

```go
// main.go — быстрая проверка до Wire
for _, arg := range os.Args[1:] {
    if arg == "-version" || arg == "--version" {
        fmt.Println(version.Version)
        os.Exit(0)
    }
}
```

## Правила изменения

- Версия обновляется вручную при релизе (формат `vMAJOR.MINOR.PATCH`)
- Префикс `v` обязателен — соответствует git-tag convention
- Не дублировать версию в Makefile, README или config.yaml
