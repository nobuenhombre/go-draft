//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"{{.ModulePath}}/src/internal/app/{{.AppName}}/cli"
	domainapp "{{.ModulePath}}/src/internal/app/{{.AppName}}/domain"
)

// initializeApp is the Wire injector entrypoint. It aggregates all ProviderSets
// and constructs the top-level application. No logic belongs here.
func initializeApp() (IApp, func(), error) {
	wire.Build(
		cli.ProviderSet,
		domainapp.ProviderSet,
		newApp,
	)
	return nil, nil, nil
}