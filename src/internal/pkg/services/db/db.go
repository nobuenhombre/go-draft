package db

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/nobuenhombre/go-draft/src/internal/pkg/locator"
	"github.com/nobuenhombre/go-draft/src/internal/pkg/services/dirs"
	"github.com/nobuenhombre/suikat/pkg/ge"
)

var (
	ErrorModulePathNotFound = errors.New("module path not found in go.mod")
)

type Service interface {
	CreateDb(dbName string) error
}

type Provider struct {
	dirs dirs.Service
}

type TemplateData struct {
	DbName     string
	ModulePath string
}

func New(dirsService dirs.Service) Service {
	return &Provider{dirs: dirsService}
}

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

	// Step 1: create directories from db template
	err = p.dirs.CreateDirs("db/", "dirs", map[string]string{"DbName": dbName})
	if err != nil {
		return ge.Pin(err)
	}

	// Step 2: process all .tpl files
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

func (p *Provider) processTemplate(tplPath, templateDir, projectRoot string, data TemplateData) error {
	relPath, err := filepath.Rel(templateDir, tplPath)
	if err != nil {
		return ge.Pin(err)
	}

	// Determine if this is a raw xo template (under sql/templates/)
	isRaw := strings.Contains(relPath, "/sql/templates/")

	content, err := os.ReadFile(tplPath)
	if err != nil {
		return ge.Pin(err)
	}

	var source []byte

	if isRaw {
		// Raw copy — no Go template processing
		source = content
	} else {
		tmpl, err := template.New(filepath.Base(tplPath)).Parse(string(content))
		if err != nil {
			return ge.Pin(err)
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, data)
		if err != nil {
			return ge.Pin(err)
		}
		source = buf.Bytes()
	}

	var outputRelPath string
	if isRaw {
		// Raw xo templates keep .tpl extension — xo reads them as .go.tpl
		outputRelPath = relPath
	} else {
		outputRelPath = strings.TrimSuffix(relPath, ".tpl")
	}
	outputRelPath = strings.ReplaceAll(outputRelPath, "{{.DbName}}", data.DbName)

	var outputPath string
	if strings.HasPrefix(outputRelPath, "_root_/") {
		outputPath = filepath.Join(projectRoot, strings.TrimPrefix(outputRelPath, "_root_/"))
	} else {
		outputPath = filepath.Join(projectRoot, "src", outputRelPath)
	}

	if _, err := os.Stat(outputPath); err == nil {
		return nil
	}

	err = os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return ge.Pin(err)
	}

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
