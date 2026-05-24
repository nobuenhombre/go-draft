# templates — Шаблоны структур проектов

## Назначение

Директория содержит YAML-шаблоны для генерации структуры Go-проектов. Используется сервисом `dirs` для создания иерархии директорий.

## Структура

```
templates/
└── dirs/
    ├── classic/     # Классическая структура Go-проекта
    │   └── config.yaml
    └── ddd/         # DDD Hexagonal Clean Architecture
        └── config.yaml
```

## Шаблоны

| Шаблон | Описание |
|--------|----------|
| `classic` | Стандартная структура: `cmd/`, `internal/app/`, `internal/pkg/`, `configs/`, `service/` |
| `ddd` | DDD Hexagonal Architecture: `layers/core/domain/`, `layers/core/application/usecases/`, `layers/adapters/` |

## Формат config.yaml

```yaml
name: <template-name>
description: <description>
variables:
    - VAR_NAME_1
    - VAR_NAME_2
directories:
    - path: path/to/dir
    - path: path/to/${VAR_NAME_1}
      permissions: "0755"
      with_git_keep: false
```

## Правила

- Директории создаются относительно рабочей директории (где запущен go-draft)
- Переменные `${VAR_NAME}` подставляются из флага `-vars`
- Если `permissions` не указан — используется 0755
- Если `with_git_keep` не указан — создаётся `.gitkeep` (default: true)