# pkg/services — Переиспользуемые сервисы

## Назначение

Директория содержит переиспользуемые сервисы инфраструктурного слоя: создание директорий (`dirs`) и парсинг переменных (`vars`).

## Структура

```
services/
├── dirs/           # Создание директорий по YAML-шаблону (см. dirs/AGENTS.md)
│   ├── config/     # Конфигурация шаблона (см. dirs/config/AGENTS.md)
│   ├── dirs.go     # Service + Provider + New() + CreateDirs()
│   └── provider.go # Wire-провайдер
│
└── vars/           # Парсинг переменных key:value (см. vars/AGENTS.md)
    ├── vars.go     # Service + Provider + New() + Parse()
    └── provider.go # Wire-провайдер
```

## Сводка

| Пакет | Go-пакет | Файлы | Назначение |
|-------|----------|-------|------------|
| `dirs/` | `dirs` | 2 (+ config 4) | Создание директорий по YAML-шаблону |
| `vars/` | `vars` | 2 | Парсинг строки `key1:val1,key2:val2` в map |

## Wire-интеграция

Сервисы входят в `InfrastructureSet`:

| Провайдер | Вход | Выход |
|-----------|------|-------|
| `ProvideDirsService` | — | `dirs.Service` |
| `ProvideVarsService` | — | `vars.Service` |

## Правила

- Каждый сервис — отдельный пакет со своим `Service` интерфейсом
- Конструктор `New()` возвращает `Service`
- Сервисы не зависят от конфигурации (stateless)
- Ошибки оборачивать через `ge.Pin(err)`