# go-draft

**Scaffolding CLI for Go projects.** Generates directory structures and Go application skeletons with Wire DI, Gin API, cron, and graceful shutdown.

Version: **v0.5.0** • [AGENTS.md](AGENTS.md)

---

## Installation

```bash
git clone <repo> && cd go-draft
make build-app
sudo make install-app   # → /usr/local/bin/go-draft
```

Or build manually:

```bash
go build -o go-draft ./src/cmd/go-draft/
```

---

## Usage

### Generate directory structure

```bash
go-draft -make=dirs -dirs=classic -vars="PROJECT_NAME:my-project"
```

Creates `cmd/`, `internal/`, `configs/`, `service/` — standard Go project layout.

### Generate a CLI application

```bash
cd /path/to/your-go-project
go-draft -make=cli -appname=my-tool
```

Generates 21 files — a complete CLI app skeleton:

```
project-root/
├── Makefile                          # deps, wire targets
├── configs/_make_/                   # Go build variables
│   ├── config/project.mk
│   ├── config/go-build.mk           # CGO_ENABLED=0, ldflags
│   └── lib/go-build/                # RAM cache, progress bar
└── src/
    ├── cmd/my-tool/
    │   ├── main.go       # entry point (panic recovery, -version)
    │   ├── app.go        # IApp + App + newApp
    │   └── wire.go       # Wire injector
    └── internal/app/my-tool/
        ├── cli/          # CLI flags + Wire provider
        ├── config/       # YAML config (fico + yaml.v3) + tests
        ├── domain/       # Business logic stub
        ├── log/          # LogFile (Open/Close/Get)
        └── version/      # const Version
```

### Generate a service application

```bash
cd /path/to/your-go-project
go-draft -make=service -appname=my-api
```

Generates 33 files — a production-ready service skeleton

```
project-root/
├── Makefile                          # deps, wire commands
├── configs/_make_/                   # Go build variables
│   ├── config/project.mk
│   ├── config/go-build.mk           # CGO_ENABLED=0, ldflags
│   └── lib/go-build/                # RAM cache, progress bar
└── src/
    ├── cmd/my-api/
    │   ├── main.go                   # + -version before Wire init
    │   ├── app.go                    # cronJob.Start(), httpServer.Run()
    │   └── wire.go                   # 6 ProviderSets
    └── internal/app/my-api/
        ├── api/server/               # Gin HTTP + graceful shutdown
        │   ├── server.go
        │   ├── server.graceful.shutdown.go
        │   ├── provider.go
        │   ├── config/
        │   └── router/
        │       ├── router.go
        │       ├── handlers/
        │       └── middlewares/
        ├── cli/                      # -runtype=init/service
        ├── config/                   # hosts + cron sections
        ├── cron-job/
        │   ├── config/
        │   └── jobs/example/         # cron.Job + Wire provider
        ├── domain/                   # + GetConfig()
        ├── log/
        └── version/
```

### Generate database scaffolding

```bash
cd /path/to/your-go-project
go-draft -make=db -dbname=my_db
```

Generates database scripts with xo code generation:

```
src/
├── internal/pkg/db/my_db/           # Go package for future models
└── scripts/xo/
    ├── xo.sh                        # shared utilities (7 files, once)
    ├── yaml.sh
    ├── postgresql.sh
    ├── backup.sh
    ├── restore.sh
    ├── create.sh
    ├── lint.sh
    └── my_db/
        ├── Makefile
        ├── xo.yaml
        ├── migrate-up.sh
        ├── migrate-down.sh
        ├── migrate-new.sh
        ├── backups/local/
        ├── backups/production/
        ├── migrations/
        ├── sql/query/many/
        ├── sql/query/one/
        ├── sql/query/uid/
        ├── sql/query/routines/
        ├── sql/query/views/
        └── sql/templates/           # 11 xo Go templates
```

Backup databases — only `{dbname}/` is added, shared scripts preserved.

### Check version

```bash
go-draft --version
# v0.5.0
```

---

## Commands

```bash
go-draft -make=dirs          # YAML-based directory template
go-draft -make=cli           # CLI application skeleton
go-draft -make=service       # Service with API + cron
go-draft -make=db            # Database scaffolding with xo scripts
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-make` | `dirs` | Command: `dirs`, `cli`, `service`, `db` |
| `-dirs` | `classic` | YAML template name: `classic`, `ddd` |
| `-appname` | `""` | Application name for `cli` / `service` |
| `-dbname` | `""` | Database name for `db` |
| `-vars` | `""` | Variables `key1:val1,key2:val2` for `dirs` |
| `-version` | — | Show version and exit |

---

## Project architecture

```
go-draft/
├── src/
│   ├── cmd/go-draft/           # Entry point + Wire DI
│   └── internal/
│       ├── app/go-draft/       # CLI, domain, version
│       └── pkg/services/
│           ├── dirs/           # Directory creation from YAML
│           ├── vars/           # key:value parsing
│           ├── app/            # Go source generation from text/template
│           └── locator/        # Template directory search
├── templates/
│   ├── dirs/classic/           # YAML template for classic layout
│   ├── dirs/ddd/               # YAML template for DDD layout
│   └── app/
│       ├── cli/                # Go templates for CLI app (21 files)
│       └── service/            # Go templates for service (33 files)
├── configs/_make_/             # Build variables
├── sandbox/                    # Acceptance test playground
└── service/deployments/        # Deployment Makefile
```

## Tech stack

| Component | Choice |
|-----------|--------|
| Language | Go 1.26 |
| DI | Google Wire |
| CLI parsing | suikat/clivar |
| Error handling | suikat/ge |
| YAML | gopkg.in/yaml.v3 |
| Go templates | text/template + go/format |
| HTTP API (service) | Gin |
| Cron (service) | robfig/cron/v3 |

---

## For developers

```bash
make wire          # Regenerate wire_gen.go
make deps          # Reinitialize go.mod
make build-app     # Build binary (service/deployments/...)
go test ./...      # Run unit + acceptance tests
go vet ./...       # Static analysis
```

### Adding a new template

1. Create a `.tpl` file in `templates/app/{cli|service}/`
2. For root-level files, place them under `_root_/`
3. Use `{{.AppName}}` and `{{.ModulePath}}` for substitution
4. Go files are automatically formatted with `go/format`

### Acceptance tests

```bash
cd sandbox
ln -s ../templates templates
go mod init github.com/test/my-app
go-draft -make=cli -appname=my-app
go build ./...
# Cleanup (preserves AGENTS.md):
cd .. && rm -rf sandbox/src sandbox/templates sandbox/go.mod sandbox/go.sum sandbox/my-app
```

---

## License

Apache 2.0
