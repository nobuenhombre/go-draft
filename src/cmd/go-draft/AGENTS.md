# go-draft (cmd) — Точка входа

## Назначение

Точка входа приложения. Содержит `main.go`, `app.go`, `wire.go` и `wire_gen.go`.

## Файлы

| Файл | Назначение |
|------|------------|
| `main.go` | Функция `main()` — panic recovery, `-version`, `initializeApp()`, `app.Run()` |
| `app.go` | `App` / `IApp` — оркестратор, делегирует `domain.Run()` |
| `wire.go` | Wire injector (`//go:build wireinject`) — агрегирует `ProviderSet` из всех пакетов |
| `wire_gen.go` | Автогенерированный код Wire — не редактировать вручную |

## main()

1. `defer recover()` — перехват паники с записью стек-трейса
2. `-version` / `--version` — быстрая проверка: выводит `version.Version` и `os.Exit(0)` без инициализации Wire
3. `initializeApp()` — создаёт все зависимости через Wire, возвращает `IApp`
4. `defer cleanup()` — гарантирует закрытие ресурсов
5. `app.Run()` — запускает оркестратор
6. При ошибках — `log.Fatalf`

## Wire ProviderSets

Каждый пакет экспортирует свой `ProviderSet` через `provider.go`. Binary `wire.go` только агрегирует:

```
cli.ProviderSet        → ProvideCLI
dirs.ProviderSet       → ProvideDirsService
vars.ProviderSet       → ProvideVarsService
domainapp.ProviderSet  → ProvideDomain
newApp                 → IApp
```

## Правила изменения

- `main.go` должен оставаться минимальным — вся логика в `internal/`
- Новые зависимости добавлять через `provider.go` в соответствующем пакете, не в `main.go`
- После изменения `provider.go` в любом пакете — запустить `make wire` для регенерации `wire_gen.go`