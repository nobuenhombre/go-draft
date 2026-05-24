package db

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/nobuenhombre/go-draft/src/internal/pkg/services/locator"
	"github.com/nobuenhombre/suikat/pkg/ge"
)

var (
	// ErrorModulePathNotFound is returned when go.mod has no module directive.
	ErrorModulePathNotFound = errors.New("module path not found in go.mod")
)

// Service defines the interface for creating database scaffolding.
type Service interface {
	// CreateDb generates database scaffolding files from templates.
	CreateDb(dbName string) error
}

// Provider implements the Service interface.
type Provider struct{}

// TemplateData holds variables for template execution.
type TemplateData struct {
	// DbName is the name of the database (e.g., "fastmail_gpt_reply").
	DbName string
	// ModulePath is the Go module path from go.mod (e.g., "github.com/org/repo").
	ModulePath string
}

// New creates a new Provider instance.
func New() Service {
	return &Provider{}
}

// CreateDb generates database scaffolding files from templates.
func (p *Provider) CreateDb(dbName string) error {
	workDir, err := os.Getwd()
	if err != nil {
		return ge.Pin(err)
	}

	projectRoot, err := findProjectRoot(workDir)
	if err != nil {
		return ge.Pin(err)
	}

	modulePath, err := readModulePath(projectRoot)
	if err != nil {
		return ge.Pin(err)
	}

	templateDir, err := locator.FindTemplateDir("db")
	if err != nil {
		return ge.Pin(err)
	}

	data := TemplateData{
		DbName:     dbName,
		ModulePath: modulePath,
	}

	err = filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".tpl") {
			return nil
		}

		return p.processTemplate(path, templateDir, projectRoot, data)
	})
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}

// processTemplate reads a single .tpl file, executes it with data, and writes the result.
func (p *Provider) processTemplate(tplPath, templateDir, projectRoot string, data TemplateData) error {
	content, err := os.ReadFile(tplPath)
	if err != nil {
		return ge.Pin(err)
	}

	tmpl, err := template.New(filepath.Base(tplPath)).Parse(string(content))
	if err != nil {
		return ge.Pin(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return ge.Pin(err)
	}

	relPath, err := filepath.Rel(templateDir, tplPath)
	if err != nil {
		return ge.Pin(err)
	}

	outputRelPath := strings.TrimSuffix(relPath, ".tpl")
	outputRelPath = strings.ReplaceAll(outputRelPath, "{{.DbName}}", data.DbName)

	// _root_ prefix → output at project root, otherwise under src/
	var outputPath string
	if strings.HasPrefix(outputRelPath, "_root_/") {
		outputPath = filepath.Join(projectRoot, strings.TrimPrefix(outputRelPath, "_root_/"))
	} else {
		outputPath = filepath.Join(projectRoot, "src", outputRelPath)
	}

	// Skip if file already exists — shared scripts should not be overwritten
	// when generating scaffolding for a second or third database.
	if _, err := os.Stat(outputPath); err == nil {
		return nil
	}

	err = os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return ge.Pin(err)
	}

	source := buf.Bytes()
	// .sh files need execute permission
	perm := os.FileMode(0644)
	if strings.HasSuffix(outputPath, ".sh") {
		perm = 0755
	}

	err = os.WriteFile(outputPath, source, perm)
	if err != nil {
		return ge.Pin(err)
	}

	return nil
}

// findProjectRoot walks up from workDir to find the directory containing go.mod.
func findProjectRoot(workDir string) (string, error) {
	dir := workDir
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ge.Pin(errors.New("project root (go.mod) not found"))
		}
		dir = parent
	}
}

// readModulePath extracts the module path from go.mod.
func readModulePath(projectRoot string) (string, error) {
	goModPath := filepath.Join(projectRoot, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", ge.Pin(err)
	}

	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			modulePath := strings.TrimSpace(strings.TrimPrefix(line, "module"))
			if modulePath != "" {
				return modulePath, nil
			}
		}
	}

	return "", ge.Pin(ErrorModulePathNotFound)
}
