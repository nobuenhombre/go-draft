package app

import (
	"log"

	"github.com/google/wire"
)

// ProviderSet exports Wire providers for the app package.
var ProviderSet = wire.NewSet(
	ProvideAppService,
)

// ProvideAppService creates the app scaffolding service.
func ProvideAppService() (Service, func(), error) {
	cleanup := func() {
		log.Println("App Service cleanup")
	}

	appService := New()

	return appService, cleanup, nil
}
