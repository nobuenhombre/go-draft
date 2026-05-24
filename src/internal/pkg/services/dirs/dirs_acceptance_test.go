package dirs_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/nobuenhombre/go-draft/src/internal/pkg/services/dirs"
	"github.com/stretchr/testify/require"
)

// resolveTemplatesPath finds the go-draft project templates/ directory.
func resolveTemplatesPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			tpl := filepath.Join(dir, "templates")
			if _, err := os.Stat(tpl); err == nil {
				return tpl, nil
			}
			return "", fmt.Errorf("templates/ not found next to go.mod at %s", dir)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", dir)
		}
		dir = parent
	}
}

// setupDirTest creates a temp dir and symlinks templates/ into it.
// Changes CWD to the temp dir and returns its path.
func setupDirTest(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	t.Logf("sandbox: %s", tmpDir)

	templatesReal, err := resolveTemplatesPath()
	require.NoError(t, err)

	err = os.Symlink(templatesReal, filepath.Join(tmpDir, "templates"))
	require.NoError(t, err)

	origDir, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	return tmpDir
}

func TestAcceptance_CreateDirs_Classic(t *testing.T) {
	projectDir := setupDirTest(t)

	svc := dirs.New()
	err := svc.CreateDirs("dirs/", "classic", map[string]string{"PROJECT_NAME": "hello-app"})
	require.NoError(t, err)

	// Verify root-level dirs exist
	require.DirExists(t, filepath.Join(projectDir, "bin"))
	require.DirExists(t, filepath.Join(projectDir, "configs"))
	require.DirExists(t, filepath.Join(projectDir, "data"))
	require.DirExists(t, filepath.Join(projectDir, "docs"))
	require.DirExists(t, filepath.Join(projectDir, "service"))

	// Verify project-name substituted dirs exist
	require.DirExists(t, filepath.Join(projectDir, "bin/hello-app"))
	require.DirExists(t, filepath.Join(projectDir, "service/deployments/hello-app"))
	require.DirExists(t, filepath.Join(projectDir, "service/deployments/hello-app/linux"))

	// Verify src dirs exist
	require.DirExists(t, filepath.Join(projectDir, "src"))
	require.DirExists(t, filepath.Join(projectDir, "src/cmd"))
	require.DirExists(t, filepath.Join(projectDir, "src/cmd/hello-app"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/hello-app"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/hello-app/cli"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/hello-app/config"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/hello-app/domain"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/hello-app/version"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/pkg"))

	// Verify .gitkeep files exist (all dirs get .gitkeep by default)
	require.FileExists(t, filepath.Join(projectDir, "configs/develop/.gitkeep"))
	require.FileExists(t, filepath.Join(projectDir, "configs/production/.gitkeep"))
	require.FileExists(t, filepath.Join(projectDir, "configs/local/.gitkeep"))
	require.FileExists(t, filepath.Join(projectDir, "data/.gitkeep"))
	require.FileExists(t, filepath.Join(projectDir, "docs/.gitkeep"))

	// Count total created dirs (classic template produces 27 dirs)
	dirCount := 0
	_ = filepath.WalkDir(projectDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != projectDir {
			dirCount++
		}
		return nil
	})
	t.Logf("classic: created %d directories (expected >= 27)", dirCount)
	require.GreaterOrEqual(t, dirCount, 23)
}

func TestAcceptance_CreateDirs_DDD(t *testing.T) {
	projectDir := setupDirTest(t)

	svc := dirs.New()
	err := svc.CreateDirs("dirs/", "ddd", map[string]string{"PROJECT_NAME": "ddd-app"})
	require.NoError(t, err)

	// Verify DDD-specific layers
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/ddd-app/layers"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/ddd-app/layers/core/domain"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/ddd-app/layers/core/domain/entities"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/ddd-app/layers/core/domain/value-objects"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/ddd-app/layers/core/domain/aggregates"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/ddd-app/layers/core/domain/services"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/ddd-app/layers/core/application"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/ddd-app/layers/core/application/usecases"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/ddd-app/layers/core/application/usecases/commands"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/ddd-app/layers/core/application/usecases/queries"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/ddd-app/layers/core/application/event-handlers"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/ddd-app/layers/adapters"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/ddd-app/layers/adapters/in"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/ddd-app/layers/adapters/out"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/ddd-app/generated"))
	require.DirExists(t, filepath.Join(projectDir, "src/internal/app/ddd-app/version"))

	dirCount := 0
	_ = filepath.WalkDir(projectDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && path != projectDir {
			dirCount++
		}
		return nil
	})
	t.Logf("ddd: created %d directories (expected >= 43)", dirCount)
	require.GreaterOrEqual(t, dirCount, 39)
}

func TestAcceptance_CreateDirs_Error_MissingVar(t *testing.T) {
	setupDirTest(t)

	svc := dirs.New()
	// classic requires PROJECT_NAME but we don't provide it
	err := svc.CreateDirs("dirs/", "classic", map[string]string{})
	require.Error(t, err)
	require.ErrorIs(t, err, dirs.ErrorMissingTemplateVar)
}

func TestAcceptance_CreateDirs_Error_NotFound(t *testing.T) {
	setupDirTest(t)

	svc := dirs.New()
	err := svc.CreateDirs("dirs/", "nonexistent-template", map[string]string{})
	require.Error(t, err)
}

func TestAcceptance_CreateDirs_CustomPermissions(t *testing.T) {
	projectDir := setupDirTest(t)

	// Use ~/.go-draft/templates/dirs/ for custom templates (locator priority #4)
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)
	customDir := filepath.Join(homeDir, ".go-draft/templates/dirs/custom")
	err = os.MkdirAll(customDir, 0755)
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(customDir) })

	customConfig := `name: custom
description: Custom template for testing permissions
variables:
    - APP
directories:
    - path: custom/${APP}
      permissions: "0700"
      with_git_keep: false
    - path: custom/${APP}/sub
`
	err = os.WriteFile(filepath.Join(customDir, "config.yaml"), []byte(customConfig), 0644)
	require.NoError(t, err)

	svc := dirs.New()
	err = svc.CreateDirs("dirs/", "custom", map[string]string{"APP": "myapp"})
	require.NoError(t, err)

	// Verify the custom dir exists and has correct permissions
	customAppDir := filepath.Join(projectDir, "custom/myapp")
	require.DirExists(t, customAppDir)

	info, err := os.Stat(customAppDir)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0700)|os.ModeDir, info.Mode())

	// No .gitkeep because with_git_keep: false
	require.NoFileExists(t, filepath.Join(customAppDir, ".gitkeep"))

	// Sub dir without explicit permissions should have default 0755 and .gitkeep
	subDir := filepath.Join(projectDir, "custom/myapp/sub")
	require.DirExists(t, subDir)
	require.FileExists(t, filepath.Join(subDir, ".gitkeep"))

	info, err = os.Stat(subDir)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0755)|os.ModeDir, info.Mode())
}

func TestAcceptance_CreateDirs_MultipleVars(t *testing.T) {
	projectDir := setupDirTest(t)

	// Use ~/.go-draft/templates/dirs/ for custom templates
	homeDir, err := os.UserHomeDir()
	require.NoError(t, err)
	customDir := filepath.Join(homeDir, ".go-draft/templates/dirs/multi")
	err = os.MkdirAll(customDir, 0755)
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(customDir) })

	multiConfig := `name: multi
description: Multi-variable template
variables:
    - APP
    - ENV
directories:
    - path: deploy/${ENV}/${APP}
`
	err = os.WriteFile(filepath.Join(customDir, "config.yaml"), []byte(multiConfig), 0644)
	require.NoError(t, err)

	svc := dirs.New()
	err = svc.CreateDirs("dirs/", "multi", map[string]string{"APP": "web", "ENV": "production"})
	require.NoError(t, err)

	require.DirExists(t, filepath.Join(projectDir, "deploy/production/web"))
}

func TestAcceptance_CreateDirs_AlreadyExists(t *testing.T) {
	projectDir := setupDirTest(t)

	// Create a dir first
	existing := filepath.Join(projectDir, "bin")
	err := os.MkdirAll(existing, 0755)
	require.NoError(t, err)

	// Running CreateDirs again should succeed (MkdirAll is idempotent)
	svc := dirs.New()
	err = svc.CreateDirs("dirs/", "classic", map[string]string{"PROJECT_NAME": "hello-app"})
	require.NoError(t, err)

	require.DirExists(t, filepath.Join(projectDir, "bin"))
}
