# vars — Парсинг переменных

## Назначение

Пакет `vars` разбирает строку переменных в формате `key1:val1,key2:val2` в `map[string]string`.

## Файлы

| Файл | Назначение |
|------|------------|
| `vars.go` | `Service` (интерфейс), `Provider` (реализация), `New()`, `Parse()` |
| `provider.go` | Wire-провайдер: `ProviderSet` + `ProvideVarsService` |

## Типы

- **`Service`** — интерфейс с методом `Parse(vars string) (map[string]string, error)`
- **`Provider`** — реализация (пустая структура, stateless)

## Поведение

`Parse("key1:val1,key2:val2")` → `{"key1": "val1", "key2": "val2"}`

Пустая строка → пустая map (без ошибки).

## Ошибки

| Ошибка | Условие |
|--------|---------|
| `ErrorInvalidKeyValuePair` | Пара не содержит `:` (например `"key1val1"`) |
| `ErrorEmptyKey` | Ключ пустой (например `":val1"`) |

## Правила

- Пакет имеет имя `vars` (без конфликтов со stdlib)
- Импортируется как `"github.com/nobuenhombre/go-draft/src/internal/pkg/services/vars"`
- Ошибки оборачивать через `ge.Pin(err)`