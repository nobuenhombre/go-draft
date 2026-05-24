package dirs

import (
	"log"

	"github.com/google/wire"
)

// ProviderSet exports Wire providers for the dirs package.
var ProviderSet = wire.NewSet(
	ProvideDirsService,
)

// ProvideDirsService creates the directory structure service.
func ProvideDirsService() (Service, func(), error) {
	cleanup := func() {
		log.Println("Dirs Service cleanup")
	}

	dirsService := New()

	return dirsService, cleanup, nil
}
