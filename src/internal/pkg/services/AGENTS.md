# pkg/services — Переиспользуемые сервисы

## Назначение

Директория содержит переиспользуемые сервисы инфраструктурного слоя: создание директорий (`dirs`), парсинг переменных (`vars`) и генерация Go-каркасов приложений (`app`).

## Структура

```
services/
├── dirs/           # Создание директорий по YAML-шаблону (см. dirs/AGENTS.md)
│   ├── config/     # Конфигурация шаблона (см. dirs/config/AGENTS.md)
│   ├── dirs.go     # Service + Provider + New() + CreateDirs()
│   └── provider.go # Wire-провайдер
│
├── vars/           # Парсинг переменных key:value (см. vars/AGENTS.md)
│   ├── vars.go     # Service + Provider + New() + Parse()
│   └── provider.go # Wire-провайдер
│
├── app/            # Генерация Go-каркасов (см. app/AGENTS.md)
│   ├── app.go      # Service + Provider + New() + CreateApp()
│   └── provider.go # Wire-провайдер
│
├── locator/        # Поиск директорий шаблонов (см. locator/AGENTS.md)
│   └── locator.go  # FindTemplateDir()
│
└── db/             # Генерация bash-скриптов БД (см. db/AGENTS.md)
    ├── db.go       # Service + Provider + New() + CreateDb()
    └── provider.go # Wire-провайдер
```

## Сводка

| Пакет | Go-пакет | Файлы | Назначение |
|-------|----------|-------|------------|
| `dirs/` | `dirs` | 2 (+ config 4) | Создание директорий по YAML-шаблону |
| `vars/` | `vars` | 2 | Парсинг строки `key1:val1,key2:val2` в map |
| `app/` | `app` | 2 | Генерация Go-файлов из text/template |

## Wire-интеграция

Сервисы входят в `InfrastructureSet`:

| Провайдер | Вход | Выход |
|-----------|------|-------|
| `ProvideDirsService` | — | `dirs.Service` |
| `ProvideVarsService` | — | `vars.Service` |
| `ProvideAppService` | — | `app.Service` |

## Правила

- Каждый сервис — отдельный пакет со своим `Service` интерфейсом
- Конструктор `New()` возвращает `Service`
- `dirs` и `vars` — stateless, `app` использует `os.Getwd()` для поиска корня проекта
- Ошибки оборачивать через `ge.Pin(err)`