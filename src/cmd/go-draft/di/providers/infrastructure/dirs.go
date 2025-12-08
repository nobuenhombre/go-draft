package infrastructure

import (
	"log"

	"github.com/nobuenhombre/go-draft/src/internal/pkg/services/dirs"
)

func ProvideDirsService() (dirs.Service, func(), error) {
	cleanup := func() {
		log.Println("Dirs Service config cleanup")
	}

	dirsService := dirs.New()

	return dirsService, cleanup, nil
}
