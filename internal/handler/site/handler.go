package site

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/handler"
	"gitlab.com/Dokuchaevvn/site-monitor/internal/handler/site/dto"
	siteusecase "gitlab.com/Dokuchaevvn/site-monitor/internal/usecase/site"
)

type getAllUseCase interface {
	Execute(ctx context.Context) ([]domain.Site, error)
}

type addUseCase interface {
	Execute(ctx context.Context, command siteusecase.AddCommand) (domain.Site, error)
}

type deleteByIDUseCase interface {
	Execute(ctx context.Context, command siteusecase.DeleteCommand) error
}

type getStatusByIDUseCase interface {
	Execute(ctx context.Context, command siteusecase.GetStatusCommand) (domain.Site, error)
}

type getHistoryByIDUseCase interface {
	Execute(ctx context.Context, command siteusecase.GetHistoryCommand) (siteusecase.GetHistoryResult, error)
}

type Handler struct {
	getAllUseCase     getAllUseCase
	addUseCase        addUseCase
	deleteUseCase     deleteByIDUseCase
	getStatusUseCase  getStatusByIDUseCase
	getHistoryUseCase getHistoryByIDUseCase
}

func NewHandler(
	getAllUseCase getAllUseCase,
	addUseCase addUseCase,
	deleteUseCase deleteByIDUseCase,
	getStatusUseCase getStatusByIDUseCase,
	getHistoryUseCase getHistoryByIDUseCase,
) *Handler {
	return &Handler{
		getAllUseCase:     getAllUseCase,
		addUseCase:        addUseCase,
		deleteUseCase:     deleteUseCase,
		getStatusUseCase:  getStatusUseCase,
		getHistoryUseCase: getHistoryUseCase,
	}
}

// GetAll godoc
// @Summary Получить список сайтов
// @Description Возвращает все сайты, добавленные в мониторинг
// @Tags sites
// @Produce json
// @Param X-Request-ID header string false "Идентификатор запроса для трассировки"
// @Success 200 {array} dto.SiteResponse
// @Failure 500 {object} handler.ErrorResponse "Внутренняя ошибка сервера"
// @Header 200 {string} X-Request-ID "Идентификатор запроса"
// @Header 500 {string} X-Request-ID "Идентификатор запроса"
// @Router /sites [get]
func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	sites, err := h.getAllUseCase.Execute(r.Context())
	if err != nil {
		mappedError := mapError(err)
		handler.WriteError(w, mappedError.Status, mappedError.Message)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(dto.ToSiteResponses(sites))
}

// Add godoc
// @Summary Добавить сайт
// @Description Создаёт новый сайт для мониторинга
// @Tags sites
// @Accept json
// @Produce json
// @Param X-Request-ID header string false "Идентификатор запроса для трассировки"
// @Param request body dto.AddSiteRequest true "Данные сайта"
// @Success 201 {object} dto.SiteResponse
// @Failure 400 {object} handler.ErrorResponse "Некорректный JSON, URL или имя сайта"
// @Failure 409 {object} handler.ErrorResponse "Сайт с таким URL уже существует"
// @Failure 500 {object} handler.ErrorResponse "Внутренняя ошибка сервера"
// @Header 201 {string} X-Request-ID "Идентификатор запроса"
// @Header 400 {string} X-Request-ID "Идентификатор запроса"
// @Header 409 {string} X-Request-ID "Идентификатор запроса"
// @Header 500 {string} X-Request-ID "Идентификатор запроса"
// @Router /sites [post]
func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	var request dto.AddSiteRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		handler.WriteError(w, http.StatusBadRequest, messageInvalidJSON)
		return
	}

	site, err := h.addUseCase.Execute(r.Context(), siteusecase.AddCommand{
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

// DeleteByID godoc
// @Summary Удалить сайт
// @Description Удаляет сайт из мониторинга по UUID
// @Tags sites
// @Produce json
// @Param X-Request-ID header string false "Идентификатор запроса для трассировки"
// @Param id path string true "UUID сайта" Format(uuid)
// @Success 204 "Сайт удалён"
// @Failure 400 {object} handler.ErrorResponse "Некорректный UUID сайта"
// @Failure 404 {object} handler.ErrorResponse "Сайт не найден"
// @Failure 500 {object} handler.ErrorResponse "Внутренняя ошибка сервера"
// @Header 204 {string} X-Request-ID "Идентификатор запроса"
// @Header 400 {string} X-Request-ID "Идентификатор запроса"
// @Header 404 {string} X-Request-ID "Идентификатор запроса"
// @Header 500 {string} X-Request-ID "Идентификатор запроса"
// @Router /sites/{id} [delete]
func (h *Handler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, messageInvalidSiteID)
		return
	}

	err = h.deleteUseCase.Execute(r.Context(), siteusecase.DeleteCommand{
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

// GetStatusByID godoc
// @Summary Получить статус сайта
// @Description Возвращает результат последней проверки сайта по UUID
// @Tags sites
// @Produce json
// @Param X-Request-ID header string false "Идентификатор запроса для трассировки"
// @Param id path string true "UUID сайта" Format(uuid)
// @Success 200 {object} dto.StatusResponse
// @Failure 400 {object} handler.ErrorResponse "Некорректный UUID сайта"
// @Failure 404 {object} handler.ErrorResponse "Сайт не найден"
// @Failure 500 {object} handler.ErrorResponse "Внутренняя ошибка сервера"
// @Header 200 {string} X-Request-ID "Идентификатор запроса"
// @Header 400 {string} X-Request-ID "Идентификатор запроса"
// @Header 404 {string} X-Request-ID "Идентификатор запроса"
// @Header 500 {string} X-Request-ID "Идентификатор запроса"
// @Router /sites/{id}/status [get]
func (h *Handler) GetStatusByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, messageInvalidSiteID)
		return
	}

	site, err := h.getStatusUseCase.Execute(r.Context(), siteusecase.GetStatusCommand{
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

// GetHistoryByID godoc
// @Summary Получить историю проверок сайта
// @Description Возвращает пагинированную историю проверок сайта по UUID
// @Tags sites
// @Produce json
// @Param X-Request-ID header string false "Идентификатор запроса для трассировки"
// @Param id path string true "UUID сайта" Format(uuid)
// @Param limit query int false "Размер страницы (по умолчанию 20, максимум 100)"
// @Param offset query int false "Смещение (по умолчанию 0)"
// @Success 200 {object} dto.CheckHistoryResponse
// @Failure 400 {object} handler.ErrorResponse "Некорректный UUID сайта или параметры пагинации"
// @Failure 404 {object} handler.ErrorResponse "Сайт не найден"
// @Failure 500 {object} handler.ErrorResponse "Внутренняя ошибка сервера"
// @Header 200 {string} X-Request-ID "Идентификатор запроса"
// @Header 400 {string} X-Request-ID "Идентификатор запроса"
// @Header 404 {string} X-Request-ID "Идентификатор запроса"
// @Header 500 {string} X-Request-ID "Идентификатор запроса"
// @Router /sites/{id}/history [get]
func (h *Handler) GetHistoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, messageInvalidSiteID)
		return
	}

	limit, err := parseOptionalInt(r.URL.Query().Get("limit"), messageInvalidLimit)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	offset, err := parseOptionalInt(r.URL.Query().Get("offset"), messageInvalidOffset)
	if err != nil {
		handler.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.getHistoryUseCase.Execute(r.Context(), siteusecase.GetHistoryCommand{
		SiteID: id,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		mappedError := mapError(err)
		handler.WriteError(w, mappedError.Status, mappedError.Message)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(dto.ToCheckHistoryResponse(
		result.Items,
		result.Total,
		result.Limit,
		result.Offset,
	))
}

func parseOptionalInt(raw, invalidMessage string) (int, error) {
	if raw == "" {
		return 0, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, errors.New(invalidMessage)
	}

	return value, nil
}
