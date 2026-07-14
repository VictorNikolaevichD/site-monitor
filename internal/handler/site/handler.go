package site

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/handler"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/handler/site/dto"
	siteusecase "gitlab.com/Dokuchaevvn/site-monitor/internal/usecase/site"
)

type getAllUseCase interface {
	Execute() []domain.Site
}

type addUseCase interface {
	Execute(command siteusecase.AddCommand) (domain.Site, error)
}

type deleteByIDUseCase interface {
	Execute(command siteusecase.DeleteCommand) error
}

type getStatusByIDUseCase interface {
	Execute(command siteusecase.GetStatusCommand) (domain.Site, error)
}

type Handler struct {
	getAllUseCase    getAllUseCase
	addUseCase       addUseCase
	deleteUseCase    deleteByIDUseCase
	getStatusUseCase getStatusByIDUseCase
}

func NewHandler(
	getAllUseCase getAllUseCase,
	addUseCase addUseCase,
	deleteUseCase deleteByIDUseCase,
	getStatusUseCase getStatusByIDUseCase,
) *Handler {
	return &Handler{
		getAllUseCase:    getAllUseCase,
		addUseCase:       addUseCase,
		deleteUseCase:    deleteUseCase,
		getStatusUseCase: getStatusUseCase,
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, _ *http.Request) {
	sites := h.getAllUseCase.Execute()

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
		URL:  request.URL,
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

func (h *Handler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, messageInvalidSiteID)
		return
	}

	err = h.deleteUseCase.Execute(siteusecase.DeleteCommand{
		ID: id,
	})
	if err != nil {
		mappedError := mapError(err)
		handler.WriteError(w, mappedError.Status, mappedError.Message)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetStatusByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, messageInvalidSiteID)
		return
	}

	site, err := h.getStatusUseCase.Execute(siteusecase.GetStatusCommand{
		ID: id,
	})
	if err != nil {
		mappedError := mapError(err)
		handler.WriteError(w, mappedError.Status, mappedError.Message)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(dto.ToStatusResponse(site))
}
