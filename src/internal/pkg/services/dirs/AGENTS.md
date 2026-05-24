# dirs — Создание директорий по YAML-шаблону

## Назначение

Пакет `dirs` создаёт иерархию директорий с настраиваемыми правами доступа по YAML-шаблону. Поддерживает подстановку переменных в путях.

## Файлы

| Файл | Назначение |
|------|------------|
| `dirs.go` | `Service` (интерфейс), `Provider` (реализация), `New()`, `CreateDirs()` |
| `provider.go` | Wire-провайдер: `ProviderSet` + `ProvideDirsService` |
| `config/` | YAML-конфигурация шаблона (см. config/AGENTS.md) |

## Типы

- **`Service`** — интерфейс с методом `CreateDirs(name string, vars map[string]string) error`
- **`Provider`** — реализация (пустая структура, stateless)

## Поведение

1. `CreateDirs(name, vars)` — найти шаблон по имени в search paths
2. Загрузить `config.yaml` из директории шаблона
3. Для каждой директории из конфига:
   - Подставить переменные `${VAR_NAME}` в путь
   - Создать директорию с указанными правами (по умолчанию 0755)
   - Создать `.gitkeep` если `with_git_keep: true` (по умолчанию true)

## Search paths

Шаблоны ищутся в порядке приоритета:

1. `/usr/local/share/go-draft/templates/dirs/{name}/`
2. `/usr/share/go-draft/templates/dirs/{name}/`
3. `/opt/go-draft/templates/dirs/{name}/`
4. `~/.go-draft/templates/dirs/{name}/`
5. `./templates/dirs/{name}/`

## Ошибки

| Ошибка | Условие |
|--------|---------|
| `ErrorCouldNotFindTemplatesDirs` | Шаблон не найден ни в одном из search paths |
| `ErrorMissingTemplateVar` | Переменная из конфига не передана в vars |

## Правила изменения

- При добавлении нового search path: обновить `searchPaths()` в `dirs.go`
- Пакет имеет имя `dirs` (не `configdirs` — это имя подпакета `config/`)
- Импортируется как `"github.com/nobuenhombre/go-draft/src/internal/pkg/services/dirs"`