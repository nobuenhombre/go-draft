# locator — Поиск директорий шаблонов

## Назначение

Пакет `locator` предоставляет единую функцию поиска директорий шаблонов в стандартных путях файловой системы. Используется сервисами `dirs` и `app` для устранения дублирования логики поиска.

## Функция

```go
func FindTemplateDir(subpath string) (string, error)
```

Ищет `templates/{subpath}` в следующем порядке:

1. `/usr/local/share/go-draft/templates/{subpath}`
2. `/usr/share/go-draft/templates/{subpath}`
3. `/opt/go-draft/templates/{subpath}`
4. `~/.go-draft/templates/{subpath}`
5. `templates/{subpath}` (относительно рабочей директории)

## Ошибки

| Ошибка | Условие |
|--------|---------|
| `ErrorTemplateDirNotFound` | Ни один путь не существует |

## Использование

```go
import "github.com/nobuenhombre/go-draft/src/internal/pkg/services/locator"

// Поиск YAML-шаблона директорий
dir, err := locator.FindTemplateDir("dirs/classic")

// Поиск Go-шаблона приложения
dir, err := locator.FindTemplateDir("app/service")
```

## Правила

- Пакет stateless — все данные получает через аргументы
- Не зависит от других пакетов проекта, только от stdlib и `suikat/ge`
- Ошибки оборачивать через `ge.Pin(err)`