# configs — Конфигурации проекта

## Назначение

Директория содержит конфигурационные файлы для разных окружений и make-библиотеки.

## Структура

```
configs/
├── develop/         # Develop-окружение (пусто — reserved)
├── local/           # Локальное окружение (пусто — reserved)
├── production/      # Продакшен-окружение (пусто — reserved)
└── _make_/          # Переменные и библиотеки для Makefile
    ├── config/
    │   ├── project.mk       # PROJECT_NAME
    │   └── go-build.mk      # Переменные сборки Go (GO_BUILD_ENV, GO_BUILD_FLAGS)
    └── lib/go-build/
        ├── cache-ram-drive.mk   # RAM-кэш для ускорения сборки
        └── progress-bar.mk      # Прогресс-бар сборки (требует gawk)
```

## Переменные сборки (go-build.mk)

| Переменная | Значение | Описание |
|------------|----------|----------|
| `GO_BUILD_ENV` | `CGO_ENABLED=0 GOOS=linux GOARCH=amd64` | Окружение статической сборки |
| `GO_BUILD_FLAGS` | `-x -ldflags="-s -w" -trimpath -tags netgo` | Флаги сборки: стриппинг, чистые пути |

## Правила

- Файлы конфигурации по окружениям заполняются вручную под каждое окружение
- Не коммитить реальные секреты в YAML-файлы
- Makefile-библиотеки подключаются через `include` из корневого и deployment Makefile