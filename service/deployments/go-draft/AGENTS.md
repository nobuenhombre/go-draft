# deployments/go-draft — Деплой приложения go-draft

## Назначение

Директория содержит скрипты сборки, установки и запуска приложения `go-draft` по платформам.

## Структура

```
go-draft/
└── linux/       # Сборка и деплой под Linux amd64 (см. linux/AGENTS.md)
    ├── .env     # Переменные окружения (пустой — reserved)
    └── Makefile # build-app, install-app, uninstall-app
```

## Поддерживаемые платформы

| Платформа | Директория | GOOS/GOARCH |
|-----------|------------|-------------|
| Linux | `linux/` | `linux/amd64` |

## Типичный процесс деплоя

```bash
cd service/deployments/go-draft/linux/

make build-app        # Сборка бинарника → bin/go-draft/linux/
make install-app      # Симлинк в /usr/local/bin/
```

## Правила

- При добавлении новой платформы — создать поддиректорию (например `darwin/`) с собственным `Makefile`
- Makefile использует `include` из `configs/_make_/` — при изменении переменных сборки проверить совместимость