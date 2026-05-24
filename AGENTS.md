# go-draft — Scaffolding CLI для Go-проектов

## Назначение

**go-draft** — CLI-утилита на Go для быстрого создания структуры Go-проектов по шаблонам. Генерирует иерархию директорий с конфигурацией из YAML-шаблонов.

## Архитектура

```
go-draft/
├── configs/                         # Конфигурации (см. configs/AGENTS.md)
│   └── _make_/                      # Переменные и библиотеки для Makefile
│       ├── config/
│       │   ├── project.mk          # PROJECT_NAME
│       │   └── go-build.mk         # Переменные сборки Go
│       └── lib/go-build/
│           ├── cache-ram-drive.mk   # RAM-кэш для сборки
│           └── progress-bar.mk      # Прогресс-бар сборки
│
├── service/                         # Сервисная инфраструктура (см. service/AGENTS.md)
│   └── deployments/go-draft/linux/  # Сборка, установка, деплой
│
├── src/                             # Исходный код (см. src/AGENTS.md)
│   ├── cmd/go-draft/                # Точка входа + Wire DI (см. cmd/go-draft/AGENTS.md)
│   └── internal/                    # Внутренние пакеты (см. internal/AGENTS.md)
│       ├── app/go-draft/            # Слой приложения (cli, domain, version, config)
│       └── pkg/services/            # Переиспользуемые сервисы (dirs, vars)
│
└── templates/                       # YAML-шаблоны структур проектов (см. templates/AGENTS.md)
    └── dirs/
        ├── classic/                 # Классическая структура Go-проекта
        └── ddd/                     # DDD Hexagonal Clean Architecture
```

## Технологический стек

| Компонент | Технология |
|-----------|------------|
| Язык | Go 1.24.3 |
| DI | Google Wire v0.7.0 |
| CLI-парсинг | suikat/clivar |
| Обработка ошибок | suikat/ge |
| Конфигурация | YAML (gopkg.in/yaml.v3) |

## Версионирование

- Единственный источник: `src/internal/app/go-draft/version/version.go`
- Формат: `vMAJOR.MINOR.PATCH`
- Флаг: `-version` / `--version` — быстрая проверка до инициализации Wire

## Conventions

- Ошибки оборачивать через `ge.Pin(err)`
- Пакеты с конфликтами имён stdlib — алиасы: `configapp`, `configdirs`
- provider.go в каждом пакете экспортирует `ProviderSet = wire.NewSet(ProvideXxx)`
- wire.go в `package main` — только импорты ProviderSet + newApp, без логики
- Именование: `Service` (интерфейс) + реализация + `New()` конструктор
- `go vet ./...` — false positive `undefined: Service` в provider.go при изолированной проверке; использовать `go vet ./...`, а не `go vet ./path/file.go`

## Команды

```bash
make wire          # Генерация wire_gen.go
make deps          # Переинициализация go.mod
make build-app     # Сборка (service/deployments/go-draft/linux/)
```

## Тестирование

```bash
go test ./...
go vet ./...
go build ./...
```

## Gotchas

- `configs/` содержит только `.gitkeep` — конфигурации заполняются вручную под окружение.
- Makefile в корне (`Makefile`) содержит общие команды; Makefile в `service/deployments/go-draft/linux/` — для сборки бинарника.