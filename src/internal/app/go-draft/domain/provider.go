package domainapp

import (
	"log"

	"github.com/google/wire"
	"github.com/nobuenhombre/go-draft/src/internal/app/go-draft/cli"
	appsvc "github.com/nobuenhombre/go-draft/src/internal/pkg/services/app"
	"github.com/nobuenhombre/go-draft/src/internal/pkg/services/dirs"
	"github.com/nobuenhombre/go-draft/src/internal/pkg/services/vars"
)

// ProviderSet exports Wire providers for the domainapp package.
var ProviderSet = wire.NewSet(
	ProvideDomain,
)

// ProvideDomain creates the domain service (business-logic orchestrator).
func ProvideDomain(cliConfig cli.Service, dirsService dirs.Service, varsService vars.Service, appService appsvc.Service) (DomainService, func(), error) {
	cleanup := func() {
		log.Println("Domain cleanup")
	}

	dom := New(cliConfig, dirsService, varsService, appService)

	return dom, cleanup, nil
}
