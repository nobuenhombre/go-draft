# scripts — Database scaffolding templates

## Назначение

Шаблоны для генерации скриптов базы данных через `-make=db`. Синхронизированы с `templates/db/`.

Сгенерированные файлы:

```
src/scripts/xo/{DbName}/
├── Makefile             # gen, backup, restore, lint, create
├── xo.yaml              # Конфигурация xo (local)
├── migrate-up.sh        # Применить миграции
├── migrate-down.sh      # Откатить миграцию
└── migrate-new.sh       # Создать новую миграцию
```