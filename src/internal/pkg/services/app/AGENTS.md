# app — Генерация Go-каркасов приложений

## Назначение

Пакет `app` генерирует Go-файлы для нового приложения из текстовых шаблонов (`text/template`). Поддерживает два типа каркасов: `cli` (консольное приложение) и `service` (системный сервис).

## Файлы

| Файл | Назначение |
|------|------------|
| `app.go` | `Service` (интерфейс), `Provider` (реализация), `New()`, `CreateApp()`, поиск шаблонов |
| `provider.go` | Wire-провайдер: `ProviderSet` + `ProvideAppService` |

## Типы

- **`Service`** — интерфейс с методом `CreateApp(appName string, appType string) error`
- **`Provider`** — реализация (stateless)
- **`TemplateData`** — `AppName` + `ModulePath` для выполнения шаблонов

## Поведение

`CreateApp("my-service", "service")`:
1. Определяет корень проекта (где `go.mod`)
2. Читает module path из `go.mod`
3. Ищет `templates/app/service/` в стандартных путях
4. Обходит все `.tpl` файлы рекурсивно
5. Для каждого:
   - Выполняет `text/template` с `{AppName, ModulePath}`
   - Если путь начинается с `_root_/` → файл в корень проекта
   - Иначе → файл в `src/{путь}`
   - Go-файлы форматируются через `go/format`

## Конвенция `_root_/`

Шаблоны в поддиректории `_root_/` создаются в корне целевого проекта:
```
templates/app/service/_root_/Makefile.tpl
  → {project}/Makefile

templates/app/service/_root_/configs/_make_/config/project.mk.tpl
  → {project}/configs/_make_/config/project.mk
```

Шаблоны вне `_root_/` создаются в `src/`:
```
templates/app/service/cmd/{{.AppName}}/main.go.tpl
  → {project}/src/cmd/{AppName}/main.go
```

## Поиск шаблонов

Шаблоны ищутся в следующем порядке (аналогично `dirs`):
1. `/usr/local/share/go-draft/templates/app/{type}/`
2. `/usr/share/go-draft/templates/app/{type}/`
3. `/opt/go-draft/templates/app/{type}/`
4. `~/.go-draft/templates/app/{type}/`
5. `./templates/app/{type}/` (относительно рабочей директории)

## Ошибки

| Ошибка | Условие |
|--------|---------|
| `ErrorCouldNotFindTemplatesApp` | Директория шаблонов не найдена |
| `ErrorModulePathNotFound` | В `go.mod` нет строки `module ...` |

## Правила

- Шаблоны — чистые Go `text/template` (не `html/template`)
- Имена файлов `.tpl` отрезаются при генерации
- `{{.AppName}}` в путях заменяется на имя приложения (используется как имя директории)
- Ошибки оборачивать через `ge.Pin(err)`