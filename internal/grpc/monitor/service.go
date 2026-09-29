package monitor

import (
	"context"

	"github.com/google/uuid"
	monitorv1 "gitlab.com/Dokuchaevvn/site-monitor/gen/go/monitor/v1"
	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
	siteusecase "gitlab.com/Dokuchaevvn/site-monitor/internal/usecase/site"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *Service) GetSite(ctx context.Context, request *monitorv1.GetSiteRequest) (*monitorv1.GetSiteResponse, error) {
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid site id")
	}

	site, err := s.getByID.Execute(ctx, siteusecase.GetByIDCommand{ID: id})
	if err != nil {
		return nil, mapError(err)
	}
	return &monitorv1.GetSiteResponse{Site: toProtoSite(site)}, nil
}

func (s *Service) CreateSite(ctx context.Context, request *monitorv1.CreateSiteRequest) (*monitorv1.CreateSiteResponse, error) {
	site, err := s.add.Execute(ctx, siteusecase.AddCommand{
		URL:  request.GetUrl(),
		Name: request.GetName(),
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &monitorv1.CreateSiteResponse{Site: toProtoSite(site)}, nil
}

func (s *Service) DeleteSite(ctx context.Context, request *monitorv1.DeleteSiteRequest) (*monitorv1.DeleteSiteResponse, error) {
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid site id")
	}

	if err := s.delete.Execute(ctx, siteusecase.DeleteCommand{ID: id}); err != nil {
		return nil, mapError(err)
	}
	return &monitorv1.DeleteSiteResponse{}, nil
}

func (s *Service) GetSiteStatus(ctx context.Context, request *monitorv1.GetSiteStatusRequest) (*monitorv1.GetSiteStatusResponse, error) {
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid site id")
	}

	site, err := s.getStatus.Execute(ctx, siteusecase.GetStatusCommand{ID: id})
	if err != nil {
		return nil, mapError(err)
	}
	return toProtoSiteStatus(site), nil
}

func (s *Service) GetSiteHistory(ctx context.Context, request *monitorv1.GetSiteHistoryRequest) (*monitorv1.GetSiteHistoryResponse, error) {
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid site id")
	}

	result, err := s.getHistory.Execute(ctx, siteusecase.GetHistoryCommand{
		SiteID: id,
		Limit:  int(request.GetLimit()),
		Offset: int(request.GetOffset()),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &monitorv1.GetSiteHistoryResponse{
		Results: toProtoCheckResults(result.Items),
		Total:   int32(result.Total),
		Limit:   int32(result.Limit),
		Offset:  int32(result.Offset),
	}, nil
}
