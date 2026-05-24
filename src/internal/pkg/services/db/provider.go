package db

import (
	"log"

	"github.com/google/wire"
)

// ProviderSet exports Wire providers for the db package.
var ProviderSet = wire.NewSet(
	ProvideDbService,
)

// ProvideDbService creates the database scaffolding service.
func ProvideDbService() (Service, func(), error) {
	cleanup := func() {
		log.Println("Db Service cleanup")
	}

	dbService := New()

	return dbService, cleanup, nil
}
