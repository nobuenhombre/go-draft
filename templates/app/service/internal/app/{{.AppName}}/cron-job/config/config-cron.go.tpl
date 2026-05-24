package configcron

import (
	configexample "{{.ModulePath}}/src/internal/app/{{.AppName}}/cron-job/jobs/example/config"
)

type CronConfig struct {
	ExampleJob configexample.ExampleJobConfig `yaml:"example_job"`
}