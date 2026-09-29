package monitor

import (
	"errors"

	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
	siteusecase "gitlab.com/Dokuchaevvn/site-monitor/internal/usecase/site"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrSiteNotFound):
		return status.Error(codes.NotFound, "site not found")
	case errors.Is(err, domain.ErrURLRequired),
		errors.Is(err, domain.ErrInvalidURL),
		errors.Is(err, domain.ErrNameRequired),
		errors.Is(err, siteusecase.ErrInvalidLimit),
		errors.Is(err, siteusecase.ErrInvalidOffset):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrSiteAlreadyExists):
		return status.Error(codes.AlreadyExists, "site already exists")
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
