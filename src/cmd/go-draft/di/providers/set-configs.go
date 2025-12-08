package providers

import (
	"github.com/google/wire"
	"github.com/nobuenhombre/go-draft/src/cmd/go-draft/di/providers/configs"
	"github.com/nobuenhombre/go-draft/src/cmd/go-draft/di/providers/domain"
	"github.com/nobuenhombre/go-draft/src/cmd/go-draft/di/providers/infrastructure"
)

var (
	// ConfigSet Конфигурация
	ConfigSet = wire.NewSet(
		configs.ProvideCLI,
	)

	InfrastructureSet = wire.NewSet(
		infrastructure.ProvideDirsService,
		infrastructure.ProvideVarsService,
	)

	DomainSet = wire.NewSet(
		domain.ProvideDomain,
	)
)
