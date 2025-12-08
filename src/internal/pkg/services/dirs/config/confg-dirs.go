package configdirs

import (
	"github.com/nobuenhombre/suikat/pkg/fico"
	"github.com/nobuenhombre/suikat/pkg/ge"
	"gopkg.in/yaml.v3"
)

// Config represents a configuration structure for managing application directories.
// It contains metadata about the configuration and defines directory structures.
type Config struct {
	Name        string      `yaml:"name"`
	Description string      `yaml:"description"`
	Variables   []string    `yaml:"variables"`
	Directories []DirConfig `yaml:"directories"`
}

// NewConfig creates a new empty Config instance.
//
// Returns:
//   - A pointer to a new Config instance
func NewConfig() *Config {
	return &Config{}
}

// Load reads and parses configuration data from a YAML file.
// The method reads the file content, unmarshals it into the Config structure,
// and returns any errors encountered during the process.
//
// Parameters:
//   - fileName: The path to the YAML configuration file
//
// Returns:
//   - An error if the file cannot be read or parsed
func (c *Config) Load(fileName string) error {
	txtConfigFile := fico.TxtFile(fileName)

	configData, err := txtConfigFile.Read()
	if err != nil {
		return ge.Pin(err)
	}

	err = yaml.Unmarshal([]byte(configData), c)
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}

// Save writes the Config structure to a YAML file.
// The method marshals the Config to YAML format and writes it to the specified file,
// returning any errors encountered during the process.
//
// Parameters:
//   - fileName: The path where the configuration should be saved
//
// Returns:
//   - An error if the data cannot be marshaled or written to the file
func (c *Config) Save(fileName string) error {
	txtConfigFile := fico.TxtFile(fileName)

	configData, err := yaml.Marshal(c)
	if err != nil {
		return ge.Pin(err)
	}

	err = txtConfigFile.Write(string(configData))
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}
