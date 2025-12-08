package infrastructure

import (
	"log"

	"github.com/nobuenhombre/go-draft/src/internal/pkg/services/vars"
)

func ProvideVarsService() (vars.Service, func(), error) {
	cleanup := func() {
		log.Println("Vars Service config cleanup")
	}

	varsService := vars.New()

	return varsService, cleanup, nil
}
