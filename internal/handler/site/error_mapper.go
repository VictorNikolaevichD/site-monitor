package site

import (
	"errors"
	"net/http"

	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
	siteusecase "gitlab.com/Dokuchaevvn/site-monitor/internal/usecase/site"
)

type ErrorMapping struct {
	Status  int
	Message string
}

func mapError(err error) ErrorMapping {
	switch {
	case errors.Is(err, domain.ErrURLRequired):
		return ErrorMapping{
			Status:  http.StatusBadRequest,
			Message: messageURLRequired,
		}
	case errors.Is(err, domain.ErrInvalidURL):
		return ErrorMapping{
			Status:  http.StatusBadRequest,
			Message: messageInvalidURL,
		}
	case errors.Is(err, domain.ErrNameRequired):
		return ErrorMapping{
			Status:  http.StatusBadRequest,
			Message: messageNameRequired,
		}
	case errors.Is(err, siteusecase.ErrSiteAlreadyExists):
		return ErrorMapping{
			Status:  http.StatusConflict,
			Message: messageSiteAlreadyExists,
		}
	default:
		return ErrorMapping{
			Status:  http.StatusInternalServerError,
			Message: messageInternalError,
		}
	}
}
