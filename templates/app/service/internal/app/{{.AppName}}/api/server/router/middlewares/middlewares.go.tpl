package middlewares

import (
	domainapp "{{.ModulePath}}/src/internal/app/{{.AppName}}/domain"
)

// HttpMiddleware holds HTTP middleware methods.
type HttpMiddleware struct {
	Domain domainapp.DomainService
}

// NewHttpMiddleware creates a new HttpMiddleware.
func NewHttpMiddleware(dom domainapp.DomainService) (mid *HttpMiddleware) {
	mid = new(HttpMiddleware)
	mid.Domain = dom
	return mid
}