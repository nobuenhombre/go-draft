package configs

import (
	"log"

	"github.com/nobuenhombre/go-draft/src/internal/app/go-draft/cli"
	"github.com/nobuenhombre/suikat/pkg/ge"
)

// ProvideCLI Провайдер CLI-конфига
func ProvideCLI() (cli.Service, func(), error) {
	cleanup := func() {
		log.Println("CLI config cleanup")
	}

	cfg, err := cli.New()
	if err != nil {
		return nil, cleanup, ge.Pin(err)
	}

	return cfg, cleanup, nil
}
