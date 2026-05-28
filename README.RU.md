# go-draft

**Скаффолдинг CLI для Go-проектов.** Генерирует структуру директорий и каркасы Go-приложений с Wire DI, Gin API, cron и graceful shutdown.

Версия: **v0.5.0** • [AGENTS.md](AGENTS.md)

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

21 файл — каркас консольного приложения с Makefile, configs/_make_/ и деплой-скриптами.

### Сгенерировать системный сервис

```bash
cd /path/to/your-go-project
go-draft -make=service -appname=my-api
```

37 файлов — Gin HTTP + graceful shutdown, cron, systemd unit-файлы.

### Сгенерировать скрипты базы данных

```bash
cd /path/to/your-go-project
go-draft -make=db -dbname=my_db
```

Создаёт scaffolding для PostgreSQL БД с xo-кодогенерацией:

```
src/
├── internal/pkg/db/my_db/           # Go-пакет для будущих моделей
└── scripts/xo/
    ├── xo.sh                        # общие утилиты (7 файлов, один раз)
    ├── yaml.sh
    ├── postgresql.sh
    ├── backup.sh
    ├── restore.sh
    ├── create.sh
    ├── lint.sh
    └── my_db/
        ├── Makefile                 # gen, backup, restore, lint, create
        ├── xo.yaml
        ├── migrate-up.sh
        ├── migrate-down.sh
        ├── migrate-new.sh
        ├── backups/local/
        ├── backups/production/
        ├── migrations/
        ├── sql/query/{many,one,uid,routines,views}/
        └── sql/templates/           # 11 xo-шаблонов Go-генерации
```

Для второй и последующих БД — только содержимое `my_db/`, общие скрипты не перезаписываются.

### Проверка версии

```bash
go-draft --version
# v0.5.0
```

---

## Команды

```bash
go-draft -make=dirs          # YAML-шаблон директорий
go-draft -make=cli           # консольное приложение
go-draft -make=service       # сервис с API + cron
go-draft -make=db            # скрипты базы данных
```

## Флаги

| Флаг | По умолчанию | Описание |
|------|-------------|----------|
| `-make` | `dirs` | Команда: `dirs`, `cli`, `service`, `db` |
| `-dirs` | `classic` | Имя YAML-шаблона: `classic`, `ddd` |
| `-appname` | `""` | Имя приложения для `cli` / `service` |
| `-dbname` | `""` | Имя базы данных для `-make=db` |
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
│           ├── app/            # Генерация Go из text/template
│           ├── locator/        # Поиск директорий шаблонов
│           └── db/             # Генерация bash-скриптов БД
├── templates/
│   ├── dirs/                   # YAML-шаблоны
│   ├── app/                    # Go-шаблоны приложений
│   └── db/                     # Bash-шаблоны БД
├── configs/_make_/             # Переменные для сборки
├── sandbox/                    # Песочница acceptance-тестов
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

1. Создать `.tpl` файл в `templates/app/{cli|service}/` или `templates/db/`
2. Для корневых файлов — положить в `_root_/`
3. Использовать `{{.AppName}}`, `{{.DbName}}`, `{{.ModulePath}}`
4. Go-файлы форматируются автоматически; `.sh` — `0755`; xo-шаблоны (в `sql/templates/`) копируются без обработки

---

## Лицензия

Apache 2.0
