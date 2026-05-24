package domainapp

import (
	"log"

	"github.com/google/wire"
	"{{.ModulePath}}/src/internal/app/{{.AppName}}/cli"
)

// ProviderSet exports Wire providers for the domainapp package.
var ProviderSet = wire.NewSet(
	ProvideDomain,
)

// ProvideDomain creates the domain service (business-logic orchestrator).
func ProvideDomain(cliConfig cli.Service) (DomainService, func(), error) {
	cleanup := func() {
		log.Println("Domain cleanup")
	}

	dom := New(cliConfig)

	return dom, cleanup, nil
}