package handler

import "net/http"

type RouteRegistrar interface {
	RegisterRoutes(mux *http.ServeMux)
}

func NewRouter(registrars ...RouteRegistrar) http.Handler {
	mux := http.NewServeMux()

	for _, registrar := range registrars {
		registrar.RegisterRoutes(mux)
	}

	return mux
}
