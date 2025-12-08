package cli

import (
	"github.com/nobuenhombre/suikat/pkg/clivar"
	"github.com/nobuenhombre/suikat/pkg/ge"
)

type Service interface {
}

// Config represents the command-line interface configuration structure.
// It holds configuration parameters for the CLI application.
type Config struct {
	Make string `cli:"make[Command to execute]:string=dirs"`
	// Dirs specifies the directory structure template for configuration files.
	Dirs string `cli:"dirs[Directory structure template]:string=classic"`
	// Vars contains variables for substitution in the format key1:value1,key2:value2.
	Vars string `cli:"vars[Variables for substitution in the format key1:value1,key2:value2]:string="`
}

// New creates a new Config instance by loading values from command-line arguments.
// It uses the clivar package to parse command-line flags into the Config struct.
//
// Returns:
//   - A pointer to the initialized Config instance
//   - An error if command-line argument parsing fails
func New() (Service, error) {
	cfg := &Config{}

	err := clivar.Load(cfg)
	if err != nil {
		return nil, ge.Pin(err)
	}

	return cfg, nil
}
