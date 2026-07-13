package site

import (
	"encoding/json"
	"net/http"

	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/handler"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/handler/site/dto"
	siteusecase "gitlab.com/Dokuchaevvn/site-monitor/internal/usecase/site"
)

type Repository interface {
	GetAll() []domain.Site
}

type AddUseCase interface {
	Execute(command siteusecase.AddCommand) (domain.Site, error)
}

type Handler struct {
	repo Repository
	addUseCase AddUseCase
}

func NewHandler(repository Repository, addUseCase AddUseCase) *Handler {
	return &Handler{
		repo: repository,
		addUseCase: addUseCase,
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
		handler.WriteError(w, http.StatusBadRequest, messageInvalidJSON)

		return
	}

	site, err := h.addUseCase.Execute(siteusecase.AddCommand{
		URL: request.URL,
		Name: request.Name,
	})
	if err != nil {
		mappedError := mapError(err)

		handler.WriteError(w, mappedError.Status, mappedError.Message)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(dto.ToSiteResponse(site))
}
