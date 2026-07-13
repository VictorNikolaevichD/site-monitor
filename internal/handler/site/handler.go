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
