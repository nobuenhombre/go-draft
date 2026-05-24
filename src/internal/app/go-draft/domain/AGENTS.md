# domainapp — Бизнес-логика (оркестратор)

## Назначение

Пакет `domainapp` — ядро приложения. Оркестрирует выполнение команд: парсинг переменных, создание директорий по YAML-шаблону.

## Структура

- **`domain-app.go`** — `AppDomain`, `DomainService`, `New()`, `Run()`, `MakeDirs()`
- **`provider.go`** — Wire-провайдер (`ProviderSet` + `ProvideDomain`)

## Ключевые типы

### AppDomain

Главная структура, хранит зависимости:

- `Cli` — `*cli.Config`
- `Dirs` — `dirs.Service`
- `Vars` — `vars.Service`

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
  │   └── default → ErrorUnknownMakeCommand
  └── return error
```

## Ошибки

| Ошибка | Условие |
|--------|---------|
| `ErrorEmptyMakeCommand` | Флаг `-make` пустой |
| `ErrorUnknownMakeCommand` | Неизвестное значение `-make` |
| `ErrorEmptyDirsTemplateName` | Флаг `-dirs` пустой |

## Требования к шаблонам

Шаблоны директорий ищутся в следующем порядке:
1. `/usr/local/share/go-draft/templates/dirs/{name}/`
2. `/usr/share/go-draft/templates/dirs/{name}/`
3. `/opt/go-draft/templates/dirs/{name}/`
4. `~/.go-draft/templates/dirs/{name}/`
5. `./templates/dirs/{name}/` (относительно рабочей директории)

## Правила

- Ошибки оборачивать через `ge.Pin(err)`
- При добавлении новой команды: добавить константу в `domain-app.go`, кейс в `switch` в `Run()`, метод в `AppDomain`