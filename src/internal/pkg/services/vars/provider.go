package vars

import (
	"log"

	"github.com/google/wire"
)

// ProviderSet exports Wire providers for the vars package.
var ProviderSet = wire.NewSet(
	ProvideVarsService,
)

// ProvideVarsService creates the variable parsing service.
func ProvideVarsService() (Service, func(), error) {
	cleanup := func() {
		log.Println("Vars Service cleanup")
	}

	varsService := New()

	return varsService, cleanup, nil
}
