# app/go-draft — Слой приложения

## Назначение

Директория содержит пакеты слоя приложения: CLI-парсинг, бизнес-логика (домен), версионирование. Все пакеты интегрируются через Wire DI.

## Структура

```
go-draft/
├── cli/           # Парсинг CLI-флагов (см. cli/AGENTS.md)
├── domain/        # Бизнес-логика (см. domain/AGENTS.md)
├── config/        # Зарезервировано для YAML-конфигурации
└── version/       # Версионирование (см. version/AGENTS.md)
```

## Пакеты

| Пакет | Go-пакет | Файлы | Назначение |
|-------|----------|-------|------------|
| `cli/` | `cli` | 2 | Парсинг CLI-флагов: `-make`, `-dirs`, `-vars` |
| `domain/` | `domainapp` | 2 | Оркестратор: `MakeDirs()`, `Run()` |
| `version/` | `version` | 1 | `const Version = "v0.1.0"` |
| `config/` | — | 0 | Reserved for future YAML config |

## Порядок инициализации (Wire)

```
1. cli.ProvideCLI()          → cli.Service             (ConfigSet)
2. dirs.ProvideDirsService() → dirs.Service             (InfrastructureSet)
3. vars.ProvideVarsService() → vars.Service             (InfrastructureSet)
4. domainapp.ProvideDomain() → domainapp.DomainService   (DomainSet)
5. newApp()                  → IApp                     (App)
```

## Поток выполнения

```
main.go → initializeApp() → app.Run()
  │
  └── domain.Run()
        ├── cli: разбор флагов (-make, -dirs, -vars)
        ├── vars: парсинг "key1:val1,key2:val2" → map
        └── dirs: создание директорий по YAML-шаблону
```

## Правила

- Именование Go-пакетов: `cli`, `domainapp` — для избежания конфликтов с stdlib
- При импорте использовать алиасы для пакетов с дефисами в пути
- Все пакеты следуют паттерну: `Service` (интерфейс) + реализация + `New()` (конструктор)
- Ошибки оборачивать через `ge.Pin(err)`