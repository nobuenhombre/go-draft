# templates — Шаблоны структур проектов и каркасов приложений

## Назначение

Директория содержит два типа шаблонов:
- **YAML-шаблоны** (`dirs/`) — для генерации структуры директорий через сервис `dirs`
- **Go-шаблоны** (`app/`) — для генерации Go-файлов каркаса приложения через сервис `app`

## Структура

```
templates/
├── dirs/                     # YAML-шаблоны для -make=dirs
│   ├── classic/              # Классическая Go-структура
│   │   └── config.yaml
│   └── ddd/                  # DDD Hexagonal Architecture
│       └── config.yaml
│
└── app/                      # Go-шаблоны для -make=cli / -make=service
    ├── cli/                  # Консольное приложение (22 .tpl файла)
    │   ├── _root_/                      # Makefile + configs/_make_/ + service/deployments/
    │   ├── cmd/{{.AppName}}/            # main.go, app.go, wire.go
    │   └── internal/app/{{.AppName}}/   # cli, config, domain, log, version
    │
    └── service/              # Системный сервис (37 .tpl файлов)
        ├── _root_/                      # Makefile + configs/_make_/ + configs/{env}/*.service + service/deployments/
        ├── cmd/{{.AppName}}/            # main.go, app.go (+ cron, API), wire.go
        ├── internal/app/{{.AppName}}/   # api/server, cli, config, cron-job, domain, log, version
        └── scripts/                     # AGENTS.md (справочно)

templates/db/                  # Bash-шаблоны для -make=db
└── scripts/xo/                # xo-утилиты и скрипты миграций
    ├── xo.sh.tpl             # Генерация Go-кода
    ├── yaml.sh.tpl           # Парсер YAML
    ├── postgresql.sh.tpl     # Строки подключения
    ├── backup.sh.tpl         # pg_dump
    ├── restore.sh.tpl        # psql restore
    ├── create.sh.tpl         # CREATE USER/DATABASE
    ├── lint.sh.tpl           # gofmt + golint
    └── {{.DbName}}/           # 5 файлов для каждой БД
        ├── Makefile.tpl
        ├── xo.yaml.tpl
        ├── migrate-up.sh.tpl
        ├── migrate-down.sh.tpl
        └── migrate-new.sh.tpl
```

## Шаблоны dirs/

| Шаблон | Описание | Файлов |
|--------|----------|--------|
| `classic` | Стандартная структура: `cmd/`, `internal/app/`, `internal/pkg/`, `configs/`, `service/` | 1 YAML |
| `ddd` | DDD Hexagonal: `layers/core/domain/`, `usecases/`, `adapters/` | 1 YAML |

## Шаблоны app/cli/ — консольное приложение (15 файлов)

```
src/
├── cmd/{AppName}/
│   ├── main.go         # точка входа (panic recovery, -version)
│   ├── app.go          # IApp + App + newApp (только domain)
│   └── wire.go         # Wire injector (4 ProviderSet)
└── internal/app/{AppName}/
    ├── cli/            # Service + Config + provider.go
    ├── config/         # YAML config (fico + yaml.v3) + тесты
    ├── domain/         # DomainService + AppDomain + provider.go
    ├── log/            # LogFile (Open/Close/Get) + provider.go
    └── version/        # const Version
```

## Шаблоны app/service/ — системный сервис (33 файла)

```
project-root/
├── Makefile                     # deps, wire — include configs/_make_/
├── configs/_make_/
│   ├── config/project.mk        # PROJECT_NAME
│   ├── config/go-build.mk       # CGO_ENABLED=0, ldflags
│   └── lib/go-build/            # cache-ram-drive, progress-bar
└── src/
    ├── cmd/{AppName}/
    │   ├── main.go              # + -version до Wire
    │   ├── app.go               # + cronJob.Start(), httpServer.Run()
    │   └── wire.go              # 6 ProviderSet
    └── internal/app/{AppName}/
        ├── api/server/          # Gin HTTP + graceful shutdown
        ├── cli/                 # + RunTypeInit, RunTypeService
        ├── config/              # + Hosts, Cron секции + тесты
        ├── cron-job/
        │   ├── config/          # CronConfig
        │   └── jobs/example/    # Job{cron.Job} + provider.go
        ├── domain/              # + GetConfig()
        ├── log/
        └── version/
```

## Переменные шаблонов app/

| Переменная | Описание | Источник |
|------------|----------|----------|
| `{{.AppName}}` | Имя приложения | Флаг `-appname` |
| `{{.ModulePath}}` | Go module path | Из `go.mod` проекта |

## Конвенция `_root_/`

Шаблоны в директории `_root_/` создаются в корне целевого проекта, а не в `src/`. Используется для Makefile и инфраструктурных конфигураций.

## Формат config.yaml (dirs/)

```yaml
name: <template-name>
description: <description>
variables:
    - VAR_NAME_1
directories:
    - path: path/to/${VAR_NAME_1}
      permissions: "0755"
      with_git_keep: false
```

## Правила

- YAML-шаблоны создают только директории (через `dirs` service)
- Go-шаблоны создают файлы с форматированием через `go/format` (через `app` service)
- Переменные `${VAR_NAME}` в YAML-шаблонах подставляются из флага `-vars`
- Переменные `{{.VarName}}` в Go-шаблонах подставляются из структуры `TemplateData`
- Если `permissions` не указан — используется 0755
- Если `with_git_keep` не указан — создаётся `.gitkeep` (default: true)
- `.tpl` расширение отрезается при генерации → `main.go.tpl` → `main.go`