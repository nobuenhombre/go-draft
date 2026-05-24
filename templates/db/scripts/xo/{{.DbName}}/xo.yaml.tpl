config:
  db:
    host: 127.0.0.1
    port: 5432
    name: {{.DbName}}
    user: {{.DbName}}
    pass: 12345
    sslmode: disable
    pool_max_conns: 10
    backups:
      path: ./backups/production/
  codegen:
    path: ../../../internal/pkg/db/{{.DbName}}/
    package: {{.DbName}}
    queries: ./sql/query/
    templates: ./sql/templates/