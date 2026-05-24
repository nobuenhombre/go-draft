# service — Сервисная инфраструктура

## Назначение

Директория содержит инфраструктуру развёртывания приложений проекта: скрипты сборки, установки и запуска.

## Структура

```
service/
└── deployments/     # Скрипты деплоя (см. deployments/AGENTS.md)
    └── go-draft/    # Деплой приложения (см. go-draft/AGENTS.md)
        └── linux/   # Linux amd64 (см. linux/AGENTS.md)
```

## Правила

- Организация: `service/deployments/{app}/{platform}/`
- Каждая платформа содержит собственный `Makefile` с командами сборки и запуска
- Makefile используют `include` из `configs/_make_/`