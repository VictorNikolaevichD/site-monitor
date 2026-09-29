package monitor

import (
	"context"

	monitorv1 "gitlab.com/Dokuchaevvn/site-monitor/gen/go/monitor/v1"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
	siteusecase "gitlab.com/Dokuchaevvn/site-monitor/internal/usecase/site"
)

type getAllUseCase interface {
	Execute(ctx context.Context) ([]domain.Site, error)
}

type getByIDUseCase interface {
	Execute(ctx context.Context, command siteusecase.GetByIDCommand) (domain.Site, error)
}

type addUseCase interface {
	Execute(ctx context.Context, command siteusecase.AddCommand) (domain.Site, error)
}

type deleteUseCase interface {
	Execute(ctx context.Context, command siteusecase.DeleteCommand) error
}

type getStatusUseCase interface {
	Execute(ctx context.Context, command siteusecase.GetStatusCommand) (domain.Site, error)
}

type getHistoryUseCase interface {
	Execute(ctx context.Context, command siteusecase.GetHistoryCommand) (siteusecase.GetHistoryResult, error)
}

type Service struct {
	monitorv1.UnimplementedMonitorServiceServer

	getAll     getAllUseCase
	getByID    getByIDUseCase
	add        addUseCase
	delete     deleteUseCase
	getStatus  getStatusUseCase
	getHistory getHistoryUseCase
}

func NewService(
	getAll getAllUseCase,
	getByID getByIDUseCase,
	add addUseCase,
	deleteUC deleteUseCase,
	getStatus getStatusUseCase,
	getHistory getHistoryUseCase,
) *Service {
	return &Service{
		getAll:     getAll,
		getByID:    getByID,
		add:        add,
		delete:     deleteUC,
		getStatus:  getStatus,
		getHistory: getHistory,
	}
}

func (s *Service) GetSites(ctx context.Context, _ *monitorv1.GetSitesRequest) (*monitorv1.GetSitesResponse, error) {
	sites, err := s.getAll.Execute(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	return &monitorv1.GetSitesResponse{Sites: toProtoSites(sites)}, nil
}
