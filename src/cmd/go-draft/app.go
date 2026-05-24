package main

import (
	domainapp "github.com/nobuenhombre/go-draft/src/internal/app/go-draft/domain"
)

// IApp is the top-level application orchestrator interface.
type IApp interface {
	Run() error
}

// App is the top-level application orchestrator.
type App struct {
	dom domainapp.DomainService
}

// Run executes the application based on CLI configuration.
func (a *App) Run() error {
	return a.dom.Run()
}

// newApp is the Wire provider for the top-level application.
func newApp(dom domainapp.DomainService) (IApp, func(), error) {
	cleanup := func() {
		// App-level cleanup if needed
	}

	return &App{dom: dom}, cleanup, nil
}
