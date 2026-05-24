# domainapp — Бизнес-логика (оркестратор)

## Назначение

Пакет `domainapp` — ядро приложения. Оркестрирует выполнение команд: парсинг переменных, создание директорий по YAML-шаблону, генерация каркасов Go-приложений.

## Структура

- **`domain-app.go`** — `AppDomain`, `DomainService`, `New()`, `Run()`, `MakeDirs()`, `MakeApp()`
- **`provider.go`** — Wire-провайдер (`ProviderSet` + `ProvideDomain`)

## Ключевые типы

### AppDomain

Главная структура, хранит зависимости:

- `Cli` — `*cli.Config`
- `Dirs` — `dirs.Service`
- `Vars` — `vars.Service`
- `App` — `appsvc.Service` (генерация Go-каркасов)

### DomainService

Интерфейс с единственным методом: `Run() error`.

## Поток выполнения

```
AppDomain.Run()
  ├── TrimSpace(makeCmd)
  ├── switch makeCmd:
  │   ├── "dirs" → MakeDirs()
  │   │   ├── vars.Parse(d.Cli.Vars) → map[string]string
  │   │   └── dirs.CreateDirs(d.Cli.Dirs, vars) → error
  │   ├── "cli" → MakeApp("cli")
  │   │   └── app.CreateApp(appName, "cli")
  │   │       └── templates/app/cli/ → 15 файлов в src/
  │   └── "service" → MakeApp("service")
  │       └── app.CreateApp(appName, "service")
  │           ├── templates/app/service/ → 27 файлов в src/
  │           └── _root_/ → Makefile + configs/_make_/ в корне
  └── return error
```

## Константы

| Константа | Значение | Команда |
|-----------|----------|---------|
| `MakeDirs` | `"dirs"` | Создание директорий по YAML-шаблону |
| `MakeCli` | `"cli"` | Генерация консольного приложения |
| `MakeService` | `"service"` | Генерация сервиса с API и cron |

## Ошибки

| Ошибка | Условие |
|--------|---------|
| `ErrorEmptyMakeCommand` | Флаг `-make` пустой |
| `ErrorUnknownMakeCommand` | Неизвестное значение `-make` |
| `ErrorEmptyDirsTemplateName` | Флаг `-dirs` пустой при `-make=dirs` |
| `ErrorEmptyAppName` | Флаг `-appname` пустой при `-make=cli` или `-make=service` |

## Требования к шаблонам

Шаблоны директорий (`-make=dirs`) ищутся в следующем порядке:
1. `/usr/local/share/go-draft/templates/dirs/{name}/`
2. `/usr/share/go-draft/templates/dirs/{name}/`
3. `/opt/go-draft/templates/dirs/{name}/`
4. `~/.go-draft/templates/dirs/{name}/`
5. `./templates/dirs/{name}/` (относительно рабочей директории)

Шаблоны приложений (`-make=cli` / `-make=service`) ищутся аналогично:
1. `/usr/local/share/go-draft/templates/app/{type}/`
2. `/usr/share/go-draft/templates/app/{type}/`
3. `/opt/go-draft/templates/app/{type}/`
4. `~/.go-draft/templates/app/{type}/`
5. `./templates/app/{type}/` (относительно рабочей директории)

## Правила

- Ошибки оборачивать через `ge.Pin(err)`
- При добавлении новой команды: добавить константу, кейс в `switch` в `Run()`, метод в `AppDomain`