package domainapp

import (
	configapp "{{.ModulePath}}/src/internal/app/{{.AppName}}/config"
	"{{.ModulePath}}/src/internal/app/{{.AppName}}/cli"
)

type DomainService interface {
	Run() error
	GetConfig() *configapp.Config
}

type AppDomain struct {
	Cli    *cli.Config
	Config *configapp.Config
}

func New(cliConfig cli.Service, appConfig configapp.Service) DomainService {
	return &AppDomain{
		Cli:    cliConfig.(*cli.Config),
		Config: appConfig.Get(),
	}
}

func (d *AppDomain) Run() error {
	// Add your domain logic here
	return nil
}

func (d *AppDomain) GetConfig() *configapp.Config {
	return d.Config
}