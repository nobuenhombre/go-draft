package locator

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/nobuenhombre/suikat/pkg/ge"
)

var (
	// ErrorTemplateDirNotFound is returned when the template directory is not found.
	ErrorTemplateDirNotFound = errors.New("template directory not found")
)

// FindTemplateDir searches for the named template subdirectory in standard locations.
// subpath is relative to templates/, e.g. "dirs/classic" or "app/service".
func FindTemplateDir(subpath string) (string, error) {
	listPaths := []string{
		"/usr/local/share/go-draft/templates/" + subpath,
		"/usr/share/go-draft/templates/" + subpath,
		"/opt/go-draft/templates/" + subpath,
		filepath.Join(os.Getenv("HOME"), ".go-draft/templates/"+subpath),
		"templates/" + subpath,
	}

	for _, path := range listPaths {
		_, err := os.Stat(path)
		if err == nil {
			return path, nil
		}
	}

	return "", ge.Pin(ErrorTemplateDirNotFound, ge.Params{"search paths": listPaths})
}
