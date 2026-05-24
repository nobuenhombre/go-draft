//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/nobuenhombre/go-draft/src/internal/app/go-draft/cli"
	domainapp "github.com/nobuenhombre/go-draft/src/internal/app/go-draft/domain"
	"github.com/nobuenhombre/go-draft/src/internal/pkg/services/dirs"
	"github.com/nobuenhombre/go-draft/src/internal/pkg/services/vars"
)

// initializeApp is the Wire injector entrypoint. It aggregates all ProviderSets
// and constructs the top-level application. No logic belongs here.
func initializeApp() (IApp, func(), error) {
	wire.Build(
		cli.ProviderSet,
		dirs.ProviderSet,
		vars.ProviderSet,
		domainapp.ProviderSet,
		newApp,
	)
	return nil, nil, nil
}
