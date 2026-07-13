package handler

import "net/http"

type RouteRegistrar interface {
	RegisterRoutes(mux *http.ServeMux)
}

func NewRouter(mux *http.ServeMux, registrars ...RouteRegistrar) {
	for _, registrar := range registrars {
		registrar.RegisterRoutes(mux)
	}
}
