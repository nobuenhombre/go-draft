# configdirs — Конфигурация YAML-шаблона директорий

## Назначение

Пакет `configdirs` загружает и предоставляет структуру YAML-конфигурации шаблона директорий.

## Файлы

| Файл | Назначение |
|------|------------|
| `confg-dirs.go` | `Config` + `NewConfig()` + `Load()` |
| `config-dirs-dir.go` | `DirConfig` + `GetPermissions()` + `IsCreateWithGitKeep()` |
| `config-dirs_test.go` | Тесты Load/Save |
| `config-dirs_test_load.yaml` | Фикстура для TestLoad |
| `config-dirs_test_save.yaml` | Результат TestSave |

## Типы

### Config

| Поле | Тип | YAML-ключ | Описание |
|------|-----|-----------|----------|
| `Name` | `string` | `name` | Имя шаблона |
| `Description` | `string` | `description` | Описание шаблона |
| `Variables` | `[]string` | `variables` | Список ожидаемых переменных подстановки |
| `Directories` | `[]DirConfig` | `directories` | Список директорий для создания |

### DirConfig

| Поле | Тип | YAML-ключ | Описание |
|------|-----|-----------|----------|
| `Path` | `string` | `path` | Путь к директории (может содержать `${VAR}`) |
| `Permissions` | `*string` | `permissions,omitempty` | Права доступа в формате "0755" (default: 0755) |
| `WithGitKeep` | `*bool` | `with_git_keep,omitempty` | Создавать `.gitkeep` (default: true) |

## Методы DirConfig

- `GetPermissions() (os.FileMode, error)` — парсит строку permissions в FileMode (формат: "0755")
- `IsCreateWithGitKeep() bool` — возвращает true если `.gitkeep` нужно создать (default: true)

## Ошибки

| Ошибка | Условие |
|--------|---------|
| `ErrorSymbolicPermissionsMustBe4Characters` | Строка permissions не 4 символа |
| `ErrorWrongOctalPermissionsFormat` | Неверный формат восьмеричных прав |

## Правила

- Пакет имеет имя `configdirs` (не `dirs`), чтобы не конфликтовать с родительским пакетом
- Импортируется как `configdirs "github.com/nobuenhombre/go-draft/src/internal/pkg/services/dirs/config"`
- Тесты используют фикстуры с тестовыми данными
- Ошибки оборачивать через `ge.Pin(err)`