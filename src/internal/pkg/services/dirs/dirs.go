package dirs

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nobuenhombre/go-draft/src/internal/pkg/locator"
	configdirs "github.com/nobuenhombre/go-draft/src/internal/pkg/services/dirs/config"
	"github.com/nobuenhombre/suikat/pkg/ge"
)

var (
	ErrorMissingTemplateVar = errors.New("template variable is missing")
)

type Service interface {
	CreateDirs(root string, name string, vars map[string]string) error
}

type Provider struct {
}

func New() Service {
	return &Provider{}
}

func (p *Provider) CreateDirs(root string, name string, vars map[string]string) error {
	path, err := locator.FindTemplateDir(root + name)
	if err != nil {
		return ge.Pin(err)
	}

	config := configdirs.NewConfig()
	err = config.Load(path + "/config.yaml")
	if err != nil {
		return ge.Pin(err)
	}

	workDir, err := os.Getwd()
	if err != nil {
		return ge.Pin(err)
	}

	for _, dirConfig := range config.Directories {
		dir := dirConfig.Path

		// Replace Variable to Values
		for _, varName := range config.Variables {
			varValue, found := vars[varName]
			if !found {
				return ge.Pin(ErrorMissingTemplateVar, ge.Params{"varName": varName})
			}

			dir = strings.ReplaceAll(dir, "${"+varName+"}", varValue)
		}

		perms, err := dirConfig.GetPermissions()
		if err != nil {
			return ge.Pin(err)
		}

		fullDir := filepath.Join(workDir, dir)

		// Create Dir
		err = os.MkdirAll(fullDir, perms)
		if err != nil {
			return ge.Pin(err)
		}

		log.Printf("Created directory '%s' with permissions: %s\n", fullDir, perms)

		// Create .gitkeep
		if dirConfig.IsCreateWithGitKeep() {
			gitkeepPath := filepath.Join(fullDir, ".gitkeep")

			err = os.WriteFile(gitkeepPath, []byte{}, os.FileMode(0644))
			if err != nil {
				return ge.Pin(err, ge.Params{".gitkeep": gitkeepPath})
			}

			log.Println("Created .gitkeep")
		}
	}

	return nil
}
