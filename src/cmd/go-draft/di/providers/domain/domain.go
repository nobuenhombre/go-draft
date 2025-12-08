package domain

import (
	"log"

	"github.com/nobuenhombre/go-draft/src/internal/app/go-draft/cli"
	domainapp "github.com/nobuenhombre/go-draft/src/internal/app/go-draft/domain"
	"github.com/nobuenhombre/go-draft/src/internal/pkg/services/dirs"
	"github.com/nobuenhombre/go-draft/src/internal/pkg/services/vars"
)

// ProvideDomain Провайдер Домена
func ProvideDomain(cliConfig cli.Service, dirsService dirs.Service, varsService vars.Service) (domainapp.DomainService, func(), error) {
	cleanup := func() {
		log.Println("Domain config cleanup")
	}

	dom := domainapp.New(cliConfig, dirsService, varsService)

	return dom, cleanup, nil
}
