# sandbox — Песочница для acceptance-тестов

## Назначение

Директория `sandbox/` предназначена для запуска acceptance-тестов go-draft. Здесь создаются временные проекты, генерируется каркас, проверяется сборка — и затем всё удаляется.

## Правила

- **Не коммитить** содержимое `sandbox/` — директория добавлена в `.gitignore`
- Единственный файл, который остаётся в `sandbox/` — `AGENTS.md` (этот)
- Все остальные файлы — временные, созданные в процессе тестирования
- После тестов удалять только сгенерированные файлы, сохраняя `AGENTS.md`

## Настройка

Шаблоны go-draft ищутся в системных путях или относительно CWD. Чтобы они были видны из sandbox, нужен симлинк:

```bash
cd sandbox
ln -s ../templates templates
```

Либо установить go-draft через `make install-app` — тогда шаблоны копируются в `/usr/local/share/go-draft/templates/`.

## Пример acceptance-теста

```bash
cd sandbox
rm -rf *                              # очистка
ln -s ../templates templates          # доступ к шаблонам

# 1. Инициализация модуля
go mod init github.com/test/my-app

# 2. Сгенерировать CLI-приложение
go-draft -make=cli -appname=my-app

# 3. Проверить, что сгенерированный код компилируется
go build ./...

# 4. Посмотреть структуру
find src -type f | sort

# 5. Очистить (сохраняя AGENTS.md)
cd .. && rm -rf sandbox/src sandbox/templates sandbox/go.mod sandbox/go.sum sandbox/my-app 2>/dev/null; echo "cleaned"

### Service-шаблон

```bash
cd sandbox
rm -rf *
ln -s ../templates templates
go mod init github.com/test/my-svc
go-draft -make=service -appname=my-svc
go build ./...
# Очистить (сохраняя AGENTS.md)
cd .. && rm -rf sandbox/src sandbox/templates sandbox/go.mod sandbox/go.sum sandbox/my-svc 2>/dev/null; echo "cleaned"

## Почему sandbox, а не TempDir

- Позволяет визуально проверить сгененрированные файлы перед удалением
- Можно отладить шаблоны без `go test` цикла
- Единое место для всех acceptance-тестов, не разбрасывая файлы по проекту