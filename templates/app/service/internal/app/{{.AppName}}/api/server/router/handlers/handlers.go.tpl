package handlers

import (
	"net/http"

	domainapp "{{.ModulePath}}/src/internal/app/{{.AppName}}/domain"
	"github.com/gin-gonic/gin"
)

// HttpHandler holds HTTP handler methods.
type HttpHandler struct {
	Domain domainapp.DomainService
}

// NewHttpHandler creates a new HttpHandler.
func NewHttpHandler(dom domainapp.DomainService) (handler *HttpHandler) {
	handler = new(HttpHandler)
	handler.Domain = dom
	return handler
}

// DefaultHandler returns a welcome message.
func (h *HttpHandler) DefaultHandler(c *gin.Context) {
	c.String(
		http.StatusOK,
		"Welcome API Server",
	)
}