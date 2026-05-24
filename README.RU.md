# go-draft

**Скаффолдинг CLI для Go-проектов.** Генерирует структуру директорий и каркасы Go-приложений с Wire DI, Gin API, cron и graceful shutdown.

Версия: **v0.2.0** • [AGENTS.md](AGENTS.md)

---

## Установка

```bash
git clone <repo> && cd go-draft
make build-app
sudo make install-app   # → /usr/local/bin/go-draft
```

Или собрать вручную:

```bash
go build -o go-draft ./src/cmd/go-draft/
```

---

## Использование

### Создать структуру директорий

```bash
go-draft -make=dirs -dirs=classic -vars="PROJECT_NAME:my-project"
```

Генерирует `cmd/`, `internal/`, `configs/`, `service/` — стандартная Go-структура.

### Сгенерировать консольное приложение

```bash
cd /path/to/your-go-project
go-draft -make=cli -appname=my-tool
```

Создаёт 21 файл — каркас консольного приложения:

```
project-root/
├── Makefile                          # deps, wire команды
├── configs/_make_/                   # Переменные сборки Go
│   ├── config/project.mk
│   ├── config/go-build.mk           # CGO_ENABLED=0, ldflags
│   └── lib/go-build/                # RAM-кэш, прогресс-бар
└── src/
├── cmd/my-tool/
│   ├── main.go       # точка входа (panic recovery, -version)
│   ├── app.go        # IApp + App + newApp
│   └── wire.go       # Wire injector
└── internal/app/my-tool/
    ├── cli/          # CLI-флаги + Wire provider
    ├── config/       # YAML config (fico + yaml.v3) + тесты
    ├── domain/       # Бизнес-логика (stub)
    ├── log/          # LogFile (Open/Close/Get)
    └── version/      # const Version
```

### Сгенерировать системный сервис

```bash
cd /path/to/your-go-project
go-draft -make=service -appname=my-api
```

Создаёт 33 файла — полный production-ready сервис:

```
project-root/
├── Makefile                          # deps, wire команды
├── configs/_make_/                   # Переменные сборки Go
│   ├── config/project.mk
│   ├── config/go-build.mk           # CGO_ENABLED=0, ldflags
│   └── lib/go-build/                # RAM-кэш, прогресс-бар
└── src/
    ├── cmd/my-api/
    │   ├── main.go                   # + -version до Wire
    │   ├── app.go                    # cronJob.Start(), httpServer.Run()
    │   └── wire.go                   # 6 ProviderSet
    └── internal/app/my-api/
        ├── api/server/               # Gin HTTP + graceful shutdown
        │   ├── server.go
        │   ├── server.graceful.shutdown.go
        │   ├── provider.go
        │   ├── config/
        │   └── router/
        │       ├── router.go
        │       ├── handlers/
        │       └── middlewares/
        ├── cli/                      # -runtype=init/service
        ├── config/                   # hosts + cron секции
        ├── cron-job/
        │   ├── config/
        │   └── jobs/example/         # cron.Job + Wire provider
        ├── domain/                   # + GetConfig()
        ├── log/
        └── version/
```

### Проверка версии

```bash
go-draft --version
# v0.2.0
```

---

## Команды

```bash
go-draft -make=dirs          # YAML-шаблон директорий
go-draft -make=cli           # консольное приложение
go-draft -make=service       # сервис с API + cron
```

## Флаги

| Флаг | По умолчанию | Описание |
|------|-------------|----------|
| `-make` | `dirs` | Команда: `dirs`, `cli`, `service` |
| `-dirs` | `classic` | Имя YAML-шаблона: `classic`, `ddd` |
| `-appname` | `""` | Имя приложения для `cli` / `service` |
| `-vars` | `""` | Переменные `key1:val1,key2:val2` для `dirs` |
| `-version` | — | Показать версию и выйти |

---

## Архитектура проекта

```
go-draft/
├── src/
│   ├── cmd/go-draft/           # Точка входа + Wire DI
│   └── internal/
│       ├── app/go-draft/       # CLI, domain, version
│       └── pkg/services/
│           ├── dirs/           # Создание директорий по YAML
│           ├── vars/           # Парсинг key:value
│           └── app/            # Генерация Go из text/template
├── templates/
│   ├── dirs/classic/           # YAML-шаблон классической структуры
│   ├── dirs/ddd/               # YAML-шаблон DDD структуры
│   └── app/
│       ├── cli/                # Go-шаблоны консольного приложения
│       └── service/            # Go-шаблоны сервиса
├── configs/_make_/             # Переменные для сборки
└── service/deployments/        # Makefile для деплоя
```

## Технологии

| Компонент | Выбор |
|-----------|-------|
| Язык | Go 1.26 |
| DI | Google Wire |
| CLI | suikat/clivar |
| Обработка ошибок | suikat/ge |
| YAML | gopkg.in/yaml.v3 |
| Шаблоны | text/template + go/format |
| HTTP API (service) | Gin |
| Cron (service) | robfig/cron/v3 |

---

## Для разработчиков

```bash
make wire          # Регенерировать wire_gen.go
make deps          # Переинициализировать go.mod
make build-app     # Собрать бинарник (service/deployments/...)
go test ./...      # Запустить тесты
go vet ./...       # Проверить код
```

### Добавление нового шаблона

1. Создать `.tpl` файл в `templates/app/{cli|service}/`
2. Если файл должен быть в корне проекта — положить в `_root_/`
3. Использовать `{{.AppName}}` и `{{.ModulePath}}` для подстановки
4. После генерации — `go/format` применяется автоматически

---

## Лицензия

Apache 2.0
