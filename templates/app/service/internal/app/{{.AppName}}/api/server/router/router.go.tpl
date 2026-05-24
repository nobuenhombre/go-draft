package router

import (
	"io"
	"os"

	domainapp "{{.ModulePath}}/src/internal/app/{{.AppName}}/domain"
	"{{.ModulePath}}/src/internal/app/{{.AppName}}/api/server/router/handlers"
	"{{.ModulePath}}/src/internal/app/{{.AppName}}/api/server/router/middlewares"
	"github.com/gin-gonic/gin"
)

// HTTPRouter wraps the Gin engine with handlers and middlewares.
type HTTPRouter struct {
	Router      *gin.Engine
	Handlers    *handlers.HttpHandler
	Middlewares *middlewares.HttpMiddleware
}

// NewHTTPRouter creates a new HTTPRouter.
func NewHTTPRouter(logFile *os.File, dom domainapp.DomainService) (router *HTTPRouter) {
	router = new(HTTPRouter)

	if logFile != nil {
		gin.DisableConsoleColor()
		gin.DefaultWriter = io.MultiWriter(logFile)
	}

	router.Handlers = handlers.NewHttpHandler(dom)
	router.Middlewares = middlewares.NewHttpMiddleware(dom)
	router.Router = gin.Default()

	router.AddRoutes()

	return
}

// AddRoutes registers all HTTP routes.
func (r *HTTPRouter) AddRoutes() {
	r.Router.GET("/", r.Handlers.DefaultHandler)
}