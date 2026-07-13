package site

import (
	"encoding/json"
	"net/http"

	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/handler/site/dto"
)

type Repository interface {
	GetAll() []domain.Site
}

type Handler struct {
	repo Repository
}

func NewHandler(repository Repository) *Handler {
	return &Handler{
		repo: repository,
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, _ *http.Request) {
	sites := h.repo.GetAll()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dto.ToSiteResponses(sites))
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	var request dto.AddSiteRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		_ = json.NewEncoder(w).Encode(dto.ErrorResponse{
			Error: messageInvalidJSON,
		})
		return
	}
}
