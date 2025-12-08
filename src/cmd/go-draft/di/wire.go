//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/nobuenhombre/go-draft/src/cmd/go-draft/di/providers"
	domainapp "github.com/nobuenhombre/go-draft/src/internal/app/go-draft/domain"
)

func InitializeApp() (domainapp.DomainService, func(), error) {
	wire.Build(
		providers.ConfigSet,
		providers.InfrastructureSet,
		providers.DomainSet,
	)

	return nil, nil, nil
}
