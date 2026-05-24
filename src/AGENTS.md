# src — Исходный код проекта

## Назначение

Корневая директория исходного кода. Содержит Go-приложение (`cmd`, `internal`) и инфраструктурные make-библиотеки (`configs/_make_/lib`).

## Структура

```
src/
├── cmd/                 # Точки входа приложений (см. cmd/AGENTS.md)
│   └── go-draft/        # main.go + Wire DI (см. cmd/go-draft/AGENTS.md)
│
└── internal/            # Внутренние пакеты (см. internal/AGENTS.md)
    ├── app/go-draft/    # Слой приложения (см. app/go-draft/AGENTS.md)
    └── pkg/services/    # Переиспользуемые сервисы (см. pkg/services/AGENTS.md)
```

## Сводка

| Директория | Пакеты | Go-файлы | Назначение |
|------------|--------|----------|------------|
| `cmd/` | 1 + DI | 4 | Точка входа, Wire DI-граф |
| `internal/app/` | 4 | 6 | CLI, domain, version |
| `internal/pkg/services/` | 2 (+ config) | 7 | Dirs, Vars сервисы |
| **Итого** | **7+** | **17** | |

## Архитектура

```
cmd/go-draft/main.go → initializeApp()
  │
  ├── app/go-draft/cli      → CLI-флаги (-make, -dirs, -vars)
  ├── pkg/services/dirs     → Создание директорий по YAML-шаблону
  ├── pkg/services/vars     → Парсинг переменных (key:value,key:value)
  │
  └── app/go-draft/domain   → Оркестратор: MakeDirs()
```

## Правила

- `internal/` — стандартный Go-механизм инкапсуляции, пакеты недоступны извне модуля
- `app/` — пакеты, специфичные для приложения, не переиспользуемые
- `pkg/` — переиспользуемые пакеты (могут использоваться другими приложениями)