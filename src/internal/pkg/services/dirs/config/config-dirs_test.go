package configdirs

import (
	"errors"
	"reflect"
	"testing"

	"github.com/nobuenhombre/suikat/pkg/fico"
)

type testConfig struct {
	fileName    string
	fileContent string
	config      *Config
	err         error
}

func TestConfigLoad(t *testing.T) {
	right644 := "0644"
	notCreateEmpty := false

	test := &testConfig{
		fileName:    "config-dirs_test_load.yaml",
		fileContent: "",
		config: &Config{
			Name:        "classic",
			Description: "A classic Go project structure",
			Variables: []string{
				"PROJECT_NAME",
			},
			Directories: []DirConfig{
				{
					Path:        "bin",
					Permissions: &right644,
					WithGitKeep: &notCreateEmpty,
				},
				{
					Path: "bin/${PROJECT_NAME}",
				},
				{
					Path: "configs",
				},
				{
					Path: "configs/develop",
				},
				{
					Path: "configs/production",
				},
				{
					Path: "configs/local",
				},
				{
					Path: "data",
				},
				{
					Path: "docs",
				},
				{
					Path: "service",
				},
				{
					Path: "service/deployments",
				},
				{
					Path: "service/deployments/${PROJECT_NAME}",
				},
				{
					Path: "service/deployments/${PROJECT_NAME}/linux",
				},
				{
					Path: "src",
				},
				{
					Path: "src/cmd",
				},
				{
					Path: "src/cmd/${PROJECT_NAME}",
				},
				{
					Path: "src/internal",
				},
				{
					Path: "src/internal/app",
				},
				{
					Path: "src/internal/app/${PROJECT_NAME}",
				},
				{
					Path: "src/internal/app/${PROJECT_NAME}/cli",
				},
				{
					Path: "src/internal/app/${PROJECT_NAME}/config",
				},
				{
					Path: "src/internal/app/${PROJECT_NAME}/domain",
				},
				{
					Path: "src/internal/pkg",
				},
			},
		},
		err: nil,
	}

	cfg := new(Config)
	err := cfg.Load(test.fileName)

	if !(reflect.DeepEqual(cfg, test.config) && errors.Is(err, test.err)) {
		t.Errorf(
			"cfg.Load(%#v),\n Expected (cfg = %#v, err = %#v),\n Actual (cfg = %#v, err = %#v).\n",
			test.fileName, test.config, test.err, cfg, err,
		)
	}
}

func TestConfigSave(t *testing.T) {
	right644 := "0644"
	notCreateEmpty := false
	test := &testConfig{
		fileName: "config-dirs_test_save.yaml",
		fileContent: "" +
			"name: classic\n" +
			"description: A classic Go project structure\n" +
			"variables:\n" +
			"    - PROJECT_NAME\n" +
			"directories:\n" +
			"    - path: bin\n" +
			"      permissions: \"0644\"\n" +
			"      with_git_keep: false\n" +
			"    - path: bin/${PROJECT_NAME}\n" +
			"    - path: configs\n" +
			"    - path: configs/develop\n" +
			"    - path: configs/production\n" +
			"    - path: configs/local\n" +
			"    - path: data\n" +
			"    - path: docs\n" +
			"    - path: service\n" +
			"    - path: service/deployments\n" +
			"    - path: service/deployments/${PROJECT_NAME}\n" +
			"    - path: service/deployments/${PROJECT_NAME}/linux\n" +
			"    - path: src\n" +
			"    - path: src/cmd\n" +
			"    - path: src/cmd/${PROJECT_NAME}\n" +
			"    - path: src/internal\n" +
			"    - path: src/internal/app\n" +
			"    - path: src/internal/app/${PROJECT_NAME}\n" +
			"    - path: src/internal/app/${PROJECT_NAME}/cli\n" +
			"    - path: src/internal/app/${PROJECT_NAME}/config\n" +
			"    - path: src/internal/app/${PROJECT_NAME}/domain\n" +
			"    - path: src/internal/pkg\n",
		config: &Config{
			Name:        "classic",
			Description: "A classic Go project structure",
			Variables: []string{
				"PROJECT_NAME",
			},
			Directories: []DirConfig{
				{
					Path:        "bin",
					Permissions: &right644,
					WithGitKeep: &notCreateEmpty,
				},
				{
					Path: "bin/${PROJECT_NAME}",
				},
				{
					Path: "configs",
				},
				{
					Path: "configs/develop",
				},
				{
					Path: "configs/production",
				},
				{
					Path: "configs/local",
				},
				{
					Path: "data",
				},
				{
					Path: "docs",
				},
				{
					Path: "service",
				},
				{
					Path: "service/deployments",
				},
				{
					Path: "service/deployments/${PROJECT_NAME}",
				},
				{
					Path: "service/deployments/${PROJECT_NAME}/linux",
				},
				{
					Path: "src",
				},
				{
					Path: "src/cmd",
				},
				{
					Path: "src/cmd/${PROJECT_NAME}",
				},
				{
					Path: "src/internal",
				},
				{
					Path: "src/internal/app",
				},
				{
					Path: "src/internal/app/${PROJECT_NAME}",
				},
				{
					Path: "src/internal/app/${PROJECT_NAME}/cli",
				},
				{
					Path: "src/internal/app/${PROJECT_NAME}/config",
				},
				{
					Path: "src/internal/app/${PROJECT_NAME}/domain",
				},
				{
					Path: "src/internal/pkg",
				},
			},
		},
		err: nil,
	}

	cfg := test.config
	err := cfg.Save(test.fileName)

	txtConfigFile := fico.TxtFile(test.fileName)
	fileContent, errReadFile := txtConfigFile.Read()

	if errReadFile != nil {
		t.Errorf(
			"txtConfigFile.Read error %#v",
			errReadFile,
		)
	}

	if !(reflect.DeepEqual(fileContent, test.fileContent) && errors.Is(err, test.err)) {
		t.Errorf(
			"cfg.Save(%#v),\n Expected (fileContent = %#v, err = %#v),\n Actual (fileContent = %#v, err = %#v).\n",
			test.fileName, test.fileContent, test.err, fileContent, err,
		)
	}
}
