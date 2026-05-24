package app_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/nobuenhombre/go-draft/src/internal/pkg/services/app"
	"github.com/stretchr/testify/require"
)

// resolveTemplatesPath finds the go-draft project templates/ directory.
// It walks up from the test binary CWD looking for go.mod.
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

// setupProject creates a temp project dir, writes go.mod, symlinks templates,
// and changes CWD to it. Returns the project dir and a cleanup func.
func setupProject(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	t.Logf("sandbox: %s", tmpDir)

	// write go.mod
	err := os.WriteFile(
		filepath.Join(tmpDir, "go.mod"),
		[]byte("module github.com/test/acceptance\n\ngo 1.26\n"),
		0644,
	)
	require.NoError(t, err)

	// symlink templates so locator finds templates/app/{cli,service}
	templatesReal, err := resolveTemplatesPath()
	require.NoError(t, err)

	err = os.Symlink(templatesReal, filepath.Join(tmpDir, "templates"))
	require.NoError(t, err)

	// change CWD to the temp project dir
	origDir, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(origDir) })

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	return tmpDir
}

// runGoBuild runs go mod tidy, wire, then go build to verify generated code compiles.
func runGoBuild(t *testing.T, dir string) {
	t.Helper()

	// 1. go mod tidy — fetch all dependencies including wire
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = dir
	tidyOutput, err := tidy.CombinedOutput()
	require.NoError(t, err, "go mod tidy failed:\n%s", string(tidyOutput))

	// 2. wire — generate wire_gen.go from wire.go (build tag wireinject)
	// Must be run on the cmd package where wire.go lives
	wireCmd := exec.Command("go", "run", "github.com/google/wire/cmd/wire@latest", "./src/cmd/...")
	wireCmd.Dir = dir
	wireOutput, err := wireCmd.CombinedOutput()
	require.NoError(t, err, "wire failed:\n%s", string(wireOutput))

	// 3. go build — verify compilation
	build := exec.Command("go", "build", "./...")
	build.Dir = dir
	output, err := build.CombinedOutput()
	require.NoError(t, err, "go build failed:\n%s", string(output))
}

func TestAcceptance_CreateApp_CLI(t *testing.T) {
	projectDir := setupProject(t)

	svc := app.New()
	err := svc.CreateApp("my-cli", "cli")
	require.NoError(t, err)

	// root-level files (Makefile, configs/_make_/...)
	require.FileExists(t, filepath.Join(projectDir, "Makefile"))
	require.FileExists(t, filepath.Join(projectDir, "configs/_make_/config/project.mk"))

	// cmd entry point
	require.DirExists(t, filepath.Join(projectDir, "src/cmd/my-cli"))
	require.FileExists(t, filepath.Join(projectDir, "src/cmd/my-cli/main.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/cmd/my-cli/app.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/cmd/my-cli/wire.go"))

	// internal packages
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-cli/cli/cli.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-cli/cli/provider.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-cli/config/config-app.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-cli/config/provider.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-cli/domain/domain-app.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-cli/domain/provider.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-cli/version/version.go"))

	// CLI must NOT generate api/server/ or cron-job/
	require.NoFileExists(t, filepath.Join(projectDir, "src/internal/app/my-cli/api/server/server.go"))
	require.NoFileExists(t, filepath.Join(projectDir, "src/internal/app/my-cli/cron-job/jobs/example/job.go"))

	// verify generated code compiles
	runGoBuild(t, projectDir)
}

func TestAcceptance_CreateApp_Service(t *testing.T) {
	projectDir := setupProject(t)

	svc := app.New()
	err := svc.CreateApp("my-svc", "service")
	require.NoError(t, err)

	// root-level files
	require.FileExists(t, filepath.Join(projectDir, "Makefile"))
	require.FileExists(t, filepath.Join(projectDir, "configs/_make_/config/go-build.mk"))

	// API server
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-svc/api/server/server.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-svc/api/server/server.graceful.shutdown.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-svc/api/server/provider.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-svc/api/server/config/config-server.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-svc/api/server/router/router.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-svc/api/server/router/handlers/handlers.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-svc/api/server/router/middlewares/middlewares.go"))

	// cron-job
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-svc/cron-job/config/config-cron.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-svc/cron-job/config/provider.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-svc/cron-job/jobs/example/config/config-example.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-svc/cron-job/jobs/example/job.go"))
	require.FileExists(t, filepath.Join(projectDir, "src/internal/app/my-svc/cron-job/jobs/example/provider.go"))

	// config with hosts + cron sections
	configContent, err := os.ReadFile(filepath.Join(projectDir, "src/internal/app/my-svc/config/config-app.go"))
	require.NoError(t, err)
	require.Contains(t, string(configContent), "HostsConfig")
	require.Contains(t, string(configContent), "Cron  configcron.CronConfig")

	// Service wire.go should include server + examplejobs ProviderSets
	wireContent, err := os.ReadFile(filepath.Join(projectDir, "src/cmd/my-svc/wire.go"))
	require.NoError(t, err)
	require.Contains(t, string(wireContent), "server.ProviderSet")
	require.Contains(t, string(wireContent), "examplejobs.ProviderSet")

	// verify generated code compiles
	runGoBuild(t, projectDir)
}

func TestAcceptance_CreateApp_Errors(t *testing.T) {
	_ = setupProject(t)

	svc := app.New()

	// nonexistent app type → error
	err := svc.CreateApp("test", "nonexistent")
	require.Error(t, err)

	// empty app name should still process templates (app service doesn't validate name)
	err = svc.CreateApp("", "cli")
	require.NoError(t, err, "empty app name should still process templates")
}

func TestAcceptance_CreateApp_ConfigTestFixture(t *testing.T) {
	projectDir := setupProject(t)

	svc := app.New()
	err := svc.CreateApp("svc-config", "service")
	require.NoError(t, err)

	// Load test fixture should exist
	loadYAML := filepath.Join(projectDir, "src/internal/app/svc-config/config/config-app_test_load.yaml")
	require.FileExists(t, loadYAML)
	content, err := os.ReadFile(loadYAML)
	require.NoError(t, err)
	require.Contains(t, string(content), "example_job")
	require.Contains(t, string(content), "hosts")
}

func TestAcceptance_CreateApp_GoFormatCompliant(t *testing.T) {
	projectDir := setupProject(t)

	svc := app.New()
	err := svc.CreateApp("my-app", "cli")
	require.NoError(t, err)

	cmd := exec.Command("gofmt", "-l", filepath.Join(projectDir, "src"))
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "gofmt failed")
	require.Empty(t, string(output), "files need gofmt:\n%s", string(output))
}
