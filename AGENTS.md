# go-draft — Scaffolding CLI для Go-проектов

## Назначение

**go-draft** — CLI-утилита на Go для быстрого создания структуры Go-проектов. Два режима:
- **`dirs`** — создание иерархии директорий по YAML-шаблону
- **`cli` / `service`** — генерация Go-файлов каркаса приложения через текстовые шаблоны (Wire DI, CLI, config, log, domain, cron-job, API)

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
│       └── pkg/services/            # Переиспользуемые сервисы (dirs, vars, app, locator)
│
├── sandbox/                         # Песочница для acceptance-тестов (см. sandbox/AGENTS.md)
│
└── templates/                       # Шаблоны для генерации (см. templates/AGENTS.md)
    ├── dirs/                        # YAML-шаблоны структур директорий
    │   ├── classic/
    │   └── ddd/
    └── app/                         # Go-шаблоны каркасов приложений
        ├── cli/                     # 22 файла — консольное приложение
        └── service/                 # 37 файлов (+ API server, cron-job, корневой Makefile, deployment Makefile, systemd units)
```

## Технологический стек

| Компонент | Технология |
|-----------|------------|
| Язык | Go 1.26.1 |
| DI | Google Wire v0.7.0 |
| CLI-парсинг | suikat/clivar |
| Обработка ошибок | suikat/ge |
| Конфигурация | YAML (gopkg.in/yaml.v3) |
| Шаблоны Go | text/template + go/format |
| HTTP API (шаблон service) | Gin (github.com/gin-gonic/gin) |
| Cron (шаблон service) | robfig/cron/v3 |

## Версионирование

- Единственный источник: `src/internal/app/go-draft/version/version.go`
- Формат: `vMAJOR.MINOR.PATCH`
- Флаг: `-version` / `--version` — быстрая проверка до инициализации Wire

## CLI-флаги

| Флаг | По умолчанию | Описание |
|------|-------------|----------|
| `-make` | `dirs` | Команда: `dirs` / `cli` / `service` |
| `-dirs` | `classic` | Имя YAML-шаблона структуры директорий |
| `-appname` | `""` | Имя приложения для `-make=cli` или `-make=service` |
| `-vars` | `""` | Переменные `key1:val1,key2:val2` для `-make=dirs` |
| `-version` | `false` | Показать версию и выйти |

## Примеры

```bash
# Создать структуру директорий
go-draft -make=dirs -dirs=classic -vars="PROJECT_NAME:my-project"

# Сгенерировать консольное приложение (15 Go-файлов)
go-draft -make=cli -appname=my-tool

# Сгенерировать сервис с API и cron (33 файла + Makefile + _make_)
go-draft -make=service -appname=my-service
```

## Conventions

- Ошибки оборачивать через `ge.Pin(err)`
- Пакеты с конфликтами имён stdlib — алиасы: `configapp`, `configdirs`, `configcron`
- provider.go в каждом пакете экспортирует `ProviderSet = wire.NewSet(ProvideXxx)`
- wire.go в `package main` — только импорты ProviderSet + newApp, без логики
- Именование: `Service` (интерфейс) + реализация + `New()` конструктор
- `go vet ./...` — false positive `undefined: Service` в provider.go при изолированной проверке; использовать `go vet ./...`, а не `go vet ./path/file.go`
- Шаблоны `_root_/` → файлы создаются в корне проекта, а не в `src/`

## Команды

```bash
make wire          # Генерация wire_gen.go
make deps          # Переинициализация go.mod
make build-app     # Сборка (service/deployments/go-draft/linux/)
```

## Тестирование

```bash
go test ./...           # Unit-тесты
go vet ./...            # Статический анализ
go build ./...          # Сборка

# Acceptance-тесты (в sandbox/ — сохраняется только AGENTS.md)
cd sandbox && go-draft -make=cli -appname=test-app && cd .. && rm -rf sandbox/src sandbox/templates sandbox/go.mod sandbox/go.sum sandbox/test-app 2>/dev/null
```

## Gotchas

- `configs/` содержит только `.gitkeep` — конфигурации заполняются вручную под окружение.
- Makefile в корне (`Makefile`) содержит общие команды; Makefile в `service/deployments/go-draft/linux/` — для сборки бинарника.
- Шаблон `classic` генерирует устаревшую `di/` директорию — после генерации перейти на per-package providers.