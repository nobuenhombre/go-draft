//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"{{.ModulePath}}/src/internal/app/{{.AppName}}/cli"
	configapp "{{.ModulePath}}/src/internal/app/{{.AppName}}/config"
	examplejobs "{{.ModulePath}}/src/internal/app/{{.AppName}}/cron-job/jobs/example"
	domainapp "{{.ModulePath}}/src/internal/app/{{.AppName}}/domain"
	logfile "{{.ModulePath}}/src/internal/app/{{.AppName}}/log"
	"{{.ModulePath}}/src/internal/app/{{.AppName}}/api/server"
)

// initializeApp is the Wire injector entrypoint. It aggregates all ProviderSets
// and constructs the top-level application. No logic belongs here.
func initializeApp() (IApp, func(), error) {
	wire.Build(
		cli.ProviderSet,
		logfile.ProviderSet,
		configapp.ProviderSet,
		domainapp.ProviderSet,
		examplejobs.ProviderSet,
		server.ProviderSet,
		newApp,
	)
	return nil, nil, nil
}