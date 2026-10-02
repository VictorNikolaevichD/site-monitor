package swagger

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "github.com/ViktorNikolaevichD/site-monitor/docs" // registers generated Swagger spec
)

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}
