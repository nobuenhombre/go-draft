package domainapp

import (
	"errors"
	"strings"

	"github.com/nobuenhombre/go-draft/src/internal/app/go-draft/cli"
	appsvc "github.com/nobuenhombre/go-draft/src/internal/pkg/services/app"
	dbsvc "github.com/nobuenhombre/go-draft/src/internal/pkg/services/db"
	"github.com/nobuenhombre/go-draft/src/internal/pkg/services/dirs"
	"github.com/nobuenhombre/go-draft/src/internal/pkg/services/vars"
	"github.com/nobuenhombre/suikat/pkg/ge"
)

var (
	ErrorEmptyMakeCommand      = errors.New("empty make command")
	ErrorUnknownMakeCommand    = errors.New("unknown make command")
	ErrorEmptyDirsTemplateName = errors.New("empty dirs template name")
	ErrorEmptyAppName          = errors.New("empty app name")
	ErrorEmptyDbName           = errors.New("empty db name")
)

const (
	MakeDirs    = "dirs"
	MakeCli     = "cli"
	MakeService = "service"
	MakeDb      = "db"
)

type DomainService interface {
	Run() error
}

type AppDomain struct {
	Cli  *cli.Config
	Dirs dirs.Service
	Vars vars.Service
	App  appsvc.Service
	Db   dbsvc.Service
}

func New(cliConfig cli.Service, dirsService dirs.Service, varsService vars.Service, appService appsvc.Service, dbService dbsvc.Service) DomainService {
	return &AppDomain{
		Cli:  cliConfig.(*cli.Config),
		Dirs: dirsService,
		Vars: varsService,
		App:  appService,
		Db:   dbService,
	}
}

func (d *AppDomain) Run() error {
	makeCmd := strings.TrimSpace(d.Cli.Make)
	if len(makeCmd) == 0 {
		return ge.Pin(ErrorEmptyMakeCommand)
	}

	switch makeCmd {
	case MakeDirs:
		err := d.MakeDirs()
		if err != nil {
			return ge.Pin(err)
		}
	case MakeCli, MakeService:
		err := d.MakeApp(makeCmd)
		if err != nil {
			return ge.Pin(err)
		}
	case MakeDb:
		err := d.MakeDb()
		if err != nil {
			return ge.Pin(err)
		}
	default:
		return ge.Pin(ErrorUnknownMakeCommand, ge.Params{"make": makeCmd})
	}

	return nil
}

func (d *AppDomain) MakeDirs() error {
	name := strings.TrimSpace(d.Cli.Dirs)
	if len(name) == 0 {
		return ge.Pin(ErrorEmptyDirsTemplateName)
	}

	vars, err := d.Vars.Parse(d.Cli.Vars)
	if err != nil {
		return ge.Pin(err)
	}

	err = d.Dirs.CreateDirs("dirs/", name, vars)
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}

func (d *AppDomain) MakeApp(appType string) error {
	appName := strings.TrimSpace(d.Cli.AppName)
	if len(appName) == 0 {
		return ge.Pin(ErrorEmptyAppName)
	}

	err := d.App.CreateApp(appName, appType)
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}

func (d *AppDomain) MakeDb() error {
	dbName := strings.TrimSpace(d.Cli.DbName)
	if len(dbName) == 0 {
		return ge.Pin(ErrorEmptyDbName)
	}

	err := d.Db.CreateDb(dbName)
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}
