# service/deployments/go-draft/linux — Сборка и деплой под Linux

## Назначение

Директория содержит `Makefile` для сборки, установки и запуска приложения `go-draft` под Linux amd64.

## Файлы

| Файл | Назначение |
|------|------------|
| `Makefile` | Команды сборки, установки и запуска |
| `.env` | Переменные окружения (зарезервировано, пока пуст) |

## Примечание по .env

`.env` подключается через `include .env` и `export $(shell sed 's/=.*//' .env)` в Makefile. Если файл пуст — `include` игнорирует пустой файл, но секция export не ломается.

## Переменные

| Переменная | Значение | Описание |
|------------|----------|----------|
| `PROJECT` | `go-draft` | Имя проекта |
| `APP_VERSION` | `v0.0.1` | Версия приложения (в Makefile — справочно, реальная версия в `version/version.go`) |
| `APP_NAME` | `go-draft` | Имя приложения |
| `PROJECT_ROOT_PATH` | `../../../..` | Относительный путь к корню проекта |
| `INSTALL_PATH` | `/usr/local/bin` | Куда устанавливается бинарник |
| `BUILD_PLATFORM` | `linux` | Целевая ОС |
| `BIN_PATH` | `bin/go-draft/linux` | Путь к собранному бинарнику |

## Команды

| Команда | Описание |
|---------|----------|
| `make build-app` | Сборка: `CGO_ENABLED=0 GOOS=linux GOARCH=amd64`, бинарник → `bin/go-draft/linux/`, проверка версии `--version` |
| `make install-app` | Установка: симлинк бинарника в `/usr/local/bin/` |
| `make uninstall-app` | Удаление симлинка из `/usr/local/bin/` |
| `make all` | `uninstall-app` → `build-app` → `install-app` |
| `make help` | Справка по командам |

## Сборка

```bash
go build -ldflags="-s -w" -o bin/go-draft/linux/go-draft -v src/cmd/go-draft/main.go
./bin/go-draft/linux/go-draft --version
```

Статическая сборка без CGO, стриппинг символов (`-s -w`), бинарник в поддиректорию платформы.

## Правила

- Перед `make install-app` необходимо выполнить `make build-app`
- `make install-app` требует `sudo`