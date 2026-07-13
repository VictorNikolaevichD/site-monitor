package handler

import "net/http"

type routeRegistrar interface {
	RegisterRoutes(mux *http.ServeMux)
}

func NewRouter(registrars ...routeRegistrar) http.Handler {
	mux := http.NewServeMux()

	for _, registrar := range registrars {
		registrar.RegisterRoutes(mux)
	}

	return mux
}
