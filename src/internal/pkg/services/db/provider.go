package db

import (
	"log"

	"github.com/google/wire"
	"github.com/nobuenhombre/go-draft/src/internal/pkg/services/dirs"
)

// ProviderSet exports Wire providers for the db package.
var ProviderSet = wire.NewSet(
	ProvideDbService,
)

// ProvideDbService creates the database scaffolding service.
func ProvideDbService(dirsService dirs.Service) (Service, func(), error) {
	cleanup := func() {
		log.Println("Db Service cleanup")
	}

	dbService := New(dirsService)

	return dbService, cleanup, nil
}
