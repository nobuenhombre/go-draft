package configdirs

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/nobuenhombre/suikat/pkg/ge"
)

var (
	// ErrorSymbolicPermissionsMustBe4Characters is returned when symbolic permissions string
	// does not have exactly 4 characters (for example: "0755").
	ErrorSymbolicPermissionsMustBe4Characters = errors.New("symbolic permissions must be 4 characters (for example: 0755)")

	// ErrorWrongOctalPermissionsFormat is returned when permissions string has incorrect octal format.
	ErrorWrongOctalPermissionsFormat = errors.New("wrong octal permissions format")
)

// DirConfig represents configuration for a directory including its path,
// permissions, and creation behavior.
type DirConfig struct {
	// Path is the filesystem path to the directory
	Path string `yaml:"path"`

	// Permissions specifies the directory permissions in octal format (e.g., "0755").
	// If nil, default permissions (0755) will be used.
	Permissions *string `yaml:"permissions,omitempty"`

	// WithGitKeep specifies whether to create a .gitkeep file in the directory.
	// If nil, defaults to true (create .gitkeep file).
	WithGitKeep *bool `yaml:"with_git_keep,omitempty"`
}

// GetPermissions returns the directory permissions as os.FileMode.
// If Permissions field is nil, returns the default value (0755).
//
// Returns:
//   - The parsed permissions as os.FileMode
//   - An error if permissions string format is invalid
func (dc *DirConfig) GetPermissions() (os.FileMode, error) {
	if dc.Permissions == nil {
		return os.FileMode(0755), nil
	}

	permStr := *dc.Permissions

	if len(permStr) != 4 {
		return 0, ge.Pin(ErrorSymbolicPermissionsMustBe4Characters)
	}

	if !strings.HasPrefix(permStr, "0") {
		return 0, ge.Pin(ErrorWrongOctalPermissionsFormat, ge.Params{"permissions": permStr})
	}

	permStrTrimmed := strings.TrimPrefix(permStr, "0")
	perm, err := strconv.ParseUint(permStrTrimmed, 8, 32)
	if err != nil {
		return 0, ge.Pin(ErrorWrongOctalPermissionsFormat, ge.Params{"permissions": permStr})
	}

	return os.FileMode(perm), nil
}

// IsCreateWithGitKeep returns whether a .gitkeep file should be created in the directory.
// If WithGitKeep field is nil, returns true (default behavior).
//
// Returns:
//   - true if .gitkeep file should be created, false otherwise
func (dc *DirConfig) IsCreateWithGitKeep() bool {
	if dc.WithGitKeep == nil {
		return true
	}

	return *dc.WithGitKeep
}
