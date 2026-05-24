package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	configserver "{{.ModulePath}}/src/internal/app/{{.AppName}}/api/server/config"
	"{{.ModulePath}}/src/internal/app/{{.AppName}}/api/server/router"
	domainapp "{{.ModulePath}}/src/internal/app/{{.AppName}}/domain"
)

// IHTTPServer defines the HTTP server interface.
type IHTTPServer interface {
	Run() error
}

// HTTPServer wraps Gin router and http.Server with graceful shutdown.
type HTTPServer struct {
	Router *router.HTTPRouter
	Server *http.Server
}

// NewHTTPServer creates a new HTTPServer.
func NewHTTPServer(config *configserver.HTTPServerConfig, logFile *os.File, dom domainapp.DomainService) (srv *HTTPServer, err error) {
	srv = new(HTTPServer)

	srv.Router = router.NewHTTPRouter(logFile, dom)

	srv.Server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Host, config.Port),
		Handler: srv.Router.Router,
	}

	return srv, nil
}

// Run starts the HTTP server and waits for a shutdown signal.
func (srv *HTTPServer) Run() error {
	go func() {
		err := srv.Server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	srv.gracefulShutDown()

	return nil
}