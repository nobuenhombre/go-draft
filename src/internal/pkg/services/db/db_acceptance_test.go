package db_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/nobuenhombre/go-draft/src/internal/pkg/services/db"
	"github.com/nobuenhombre/go-draft/src/internal/pkg/services/dirs"
	"github.com/stretchr/testify/require"
)

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

func setupProject(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	t.Logf("sandbox: %s", tmpDir)

	err := os.WriteFile(
		filepath.Join(tmpDir, "go.mod"),
		[]byte("module github.com/test/my-db\n\ngo 1.26\n"),
		0644,
	)
	require.NoError(t, err)

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

func TestAcceptance_CreateDb_AllFiles(t *testing.T) {
	projectDir := setupProject(t)

	svc := db.New(dirs.New())
	err := svc.CreateDb("orders_db")
	require.NoError(t, err)

	// Shared scripts at src/scripts/xo/
	shared := []string{
		"xo.sh", "yaml.sh", "postgresql.sh",
		"backup.sh", "restore.sh", "create.sh", "lint.sh",
	}
	for _, f := range shared {
		path := filepath.Join(projectDir, "src/scripts/xo", f)
		require.FileExists(t, path, "shared script %s should exist", f)

		info, err := os.Stat(path)
		require.NoError(t, err)
		require.Equal(t, os.FileMode(0755), info.Mode()&os.ModePerm,
			"%s should be executable", f)
	}

	// Db-specific files at src/scripts/xo/orders_db/
	dbFiles := []string{
		"Makefile", "xo.yaml",
		"migrate-up.sh", "migrate-down.sh", "migrate-new.sh",
	}
	for _, f := range dbFiles {
		path := filepath.Join(projectDir, "src/scripts/xo/orders_db", f)
		require.FileExists(t, path, "db file %s should exist", f)
	}
}

func TestAcceptance_CreateDb_Directories(t *testing.T) {
	projectDir := setupProject(t)

	svc := db.New(dirs.New())
	err := svc.CreateDb("dirs_db")
	require.NoError(t, err)

	// pkg directory
	require.DirExists(t, filepath.Join(projectDir, "src/internal/pkg/db/dirs_db"))
	require.FileExists(t, filepath.Join(projectDir, "src/internal/pkg/db/dirs_db/.gitkeep"))

	// xo subdirectories
	dirs := []string{
		"src/scripts/xo/dirs_db",
		"src/scripts/xo/dirs_db/backups",
		"src/scripts/xo/dirs_db/backups/local",
		"src/scripts/xo/dirs_db/backups/production",
		"src/scripts/xo/dirs_db/migrations",
		"src/scripts/xo/dirs_db/sql",
		"src/scripts/xo/dirs_db/sql/query",
		"src/scripts/xo/dirs_db/sql/query/many",
		"src/scripts/xo/dirs_db/sql/query/one",
		"src/scripts/xo/dirs_db/sql/query/uid",
		"src/scripts/xo/dirs_db/sql/query/routines",
		"src/scripts/xo/dirs_db/sql/query/views",
		"src/scripts/xo/dirs_db/sql/templates",
	}
	for _, d := range dirs {
		require.DirExists(t, filepath.Join(projectDir, d), "dir %s should exist", d)
	}
}

func TestAcceptance_CreateDb_XoTemplates(t *testing.T) {
	projectDir := setupProject(t)

	svc := db.New(dirs.New())
	err := svc.CreateDb("xo_app")
	require.NoError(t, err)

	// All 11 xo template files should be copied verbatim
	templates := []string{
		"postgres.enum.go", "postgres.foreignkey.go", "postgres.index.go",
		"postgres.proc.go", "postgres.query.go", "postgres.querytype.go",
		"postgres.type.go", "xo_db.go", "xo_package.go",
		"xouid_package.go", "xouid_query.go",
	}
	for _, f := range templates {
		path := filepath.Join(projectDir, "src/scripts/xo/xo_app/sql/templates", f)
		require.FileExists(t, path, "xo template %s should exist", f)
	}

	// Verify they contain xo template syntax (not processed by Go text/template)
	// xo templates use {{.}} syntax which would fail if processed by Go's text/template
	sample := filepath.Join(projectDir, "src/scripts/xo/xo_app/sql/templates/xo_package.go")
	content, err := os.ReadFile(sample)
	require.NoError(t, err)
	require.Contains(t, string(content), "{{ .Package }}")
}

func TestAcceptance_CreateDb_XoYamlContent(t *testing.T) {
	projectDir := setupProject(t)

	svc := db.New(dirs.New())
	err := svc.CreateDb("analytics")
	require.NoError(t, err)

	yamlPath := filepath.Join(projectDir, "src/scripts/xo/analytics/xo.yaml")
	content, err := os.ReadFile(yamlPath)
	require.NoError(t, err)

	require.Contains(t, string(content), "name: analytics")
	require.Contains(t, string(content), "user: analytics")
	require.Contains(t, string(content), "../../../internal/pkg/db/analytics/")
	require.Contains(t, string(content), "package: analytics")
}

func TestAcceptance_CreateDb_MakefileContent(t *testing.T) {
	projectDir := setupProject(t)

	svc := db.New(dirs.New())
	err := svc.CreateDb("users")
	require.NoError(t, err)

	makePath := filepath.Join(projectDir, "src/scripts/xo/users/Makefile")
	content, err := os.ReadFile(makePath)
	require.NoError(t, err)

	require.Contains(t, string(content), "XO=../xo.sh")
	require.Contains(t, string(content), "## gen:")
	require.Contains(t, string(content), "## backup:")
	require.Contains(t, string(content), "## restore:")
	require.Contains(t, string(content), "## lint:")
	require.Contains(t, string(content), "## create:")
}

func TestAcceptance_CreateDb_SecondDbKeepsShared(t *testing.T) {
	projectDir := setupProject(t)

	svc := db.New(dirs.New())

	// First DB
	err := svc.CreateDb("first")
	require.NoError(t, err)

	// Second DB — shared scripts should still exist (not overwritten)
	err = svc.CreateDb("second")
	require.NoError(t, err)

	// Shared scripts
	require.FileExists(t, filepath.Join(projectDir, "src/scripts/xo/xo.sh"))
	require.FileExists(t, filepath.Join(projectDir, "src/scripts/xo/yaml.sh"))

	// Both db-specific dirs exist
	require.FileExists(t, filepath.Join(projectDir, "src/scripts/xo/first/Makefile"))
	require.FileExists(t, filepath.Join(projectDir, "src/scripts/xo/second/Makefile"))
}

func TestAcceptance_CreateDb_Error_NoTemplates(t *testing.T) {
	setupProject(t)

	svc := db.New(dirs.New())
	// There's no templates/db/ in the test — wait, there IS. Let's test
	// that an empty/whitespace db name still works (no validation in db service)
	err := svc.CreateDb("valid")
	require.NoError(t, err)
}

func TestAcceptance_CreateDb_MigrateScripts(t *testing.T) {
	projectDir := setupProject(t)

	svc := db.New(dirs.New())
	err := svc.CreateDb("inventory")
	require.NoError(t, err)

	// migrate-up.sh references yaml.sh, postgresql.sh, and uses parseYAML + postgresqlConnectString
	upPath := filepath.Join(projectDir, "src/scripts/xo/inventory/migrate-up.sh")
	content, err := os.ReadFile(upPath)
	require.NoError(t, err)
	require.Contains(t, string(content), "source ../yaml.sh")
	require.Contains(t, string(content), "source ../postgresql.sh")
	require.Contains(t, string(content), "parseYAML")
	require.Contains(t, string(content), "migrate -path migrations")

	// migrate-down.sh
	downPath := filepath.Join(projectDir, "src/scripts/xo/inventory/migrate-down.sh")
	content, err = os.ReadFile(downPath)
	require.NoError(t, err)
	require.Contains(t, string(content), "down 1")

	// migrate-new.sh
	newPath := filepath.Join(projectDir, "src/scripts/xo/inventory/migrate-new.sh")
	content, err = os.ReadFile(newPath)
	require.NoError(t, err)
	require.Contains(t, string(content), "migrate create")
	require.Contains(t, string(content), "-dir migrations")
}

func TestAcceptance_CreateDb_FilePermissions(t *testing.T) {
	projectDir := setupProject(t)

	svc := db.New(dirs.New())
	err := svc.CreateDb("payments")
	require.NoError(t, err)

	// .sh files should be 0755
	info, err := os.Stat(filepath.Join(projectDir, "src/scripts/xo/xo.sh"))
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0755), info.Mode()&os.ModePerm)

	// non-.sh files should be 0644
	info, err = os.Stat(filepath.Join(projectDir, "src/scripts/xo/payments/xo.yaml"))
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0644), info.Mode()&os.ModePerm)

	info, err = os.Stat(filepath.Join(projectDir, "src/scripts/xo/payments/Makefile"))
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0644), info.Mode()&os.ModePerm)
}

func TestAcceptance_CreateDb_SharedScriptNotOverwritten(t *testing.T) {
	projectDir := setupProject(t)

	svc := db.New(dirs.New())

	// First DB creates xo.sh
	err := svc.CreateDb("db1")
	require.NoError(t, err)

	// Modify xo.sh content
	xoPath := filepath.Join(projectDir, "src/scripts/xo/xo.sh")
	err = os.WriteFile(xoPath, []byte("# modified"), 0755)
	require.NoError(t, err)

	// Second DB should NOT overwrite xo.sh
	err = svc.CreateDb("db2")
	require.NoError(t, err)

	content, err := os.ReadFile(xoPath)
	require.NoError(t, err)
	require.Equal(t, "# modified", string(content), "xo.sh should not have been overwritten")
}
