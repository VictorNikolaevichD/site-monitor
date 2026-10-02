package client

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	monitorv1 "gitlab.com/Dokuchaevvn/site-monitor/gen/go/monitor/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type MonitorClient struct {
	raw     monitorv1.MonitorServiceClient
	timeout time.Duration
	logger  *slog.Logger
}

func NewMonitorClient(
	addr string, timeout time.Duration, logger *slog.Logger,
) (*MonitorClient, *grpc.ClientConn, error) {
	ka := keepalive.ClientParameters{
		Time:                30 * time.Second,
		Timeout:             10 * time.Second,
		PermitWithoutStream: true,
	}

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(ka),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff:           backoff.DefaultConfig,
			MinConnectTimeout: 5 * time.Second,
		}),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("dial monitor grpc %s: %w", addr, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	conn.Connect()

	for {
		state := conn.GetState()
		if state == connectivity.Ready {
			break
		}
		if !conn.WaitForStateChange(ctx, state) {
			_ = conn.Close()
			return nil, nil, fmt.Errorf("monitor grpc %s not ready: %w", addr, ctx.Err())
		}
	}

	raw := monitorv1.NewMonitorServiceClient(conn)
	return &MonitorClient{
		raw:     raw,
		timeout: timeout,
		logger:  logger,
	}, conn, nil
}

func (c *MonitorClient) GetSite(ctx context.Context, id string) (*monitorv1.GetSiteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.raw.GetSite(ctx, &monitorv1.GetSiteRequest{Id: id})
	if err != nil {
		c.logger.Error("grpc GetSite failed", "error", err)
		return nil, err
	}
	return resp, nil
}

func (c *MonitorClient) GetSites(ctx context.Context) (*monitorv1.GetSitesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.raw.GetSites(ctx, &monitorv1.GetSitesRequest{})
	if err != nil {
		c.logger.Error("grpc GetSites failed", "error", err)
		return nil, err
	}
	return resp, nil
}

func (c *MonitorClient) CreateSite(ctx context.Context, url, name string) (*monitorv1.CreateSiteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.raw.CreateSite(ctx, &monitorv1.CreateSiteRequest{
		Url:  url,
		Name: name,
	})
	if err != nil {
		c.logger.Error("grpc CreateSite failed", "error", err)
		return nil, err
	}
	return resp, nil
}

func (c *MonitorClient) DeleteSite(ctx context.Context, id string) (*monitorv1.DeleteSiteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.raw.DeleteSite(ctx, &monitorv1.DeleteSiteRequest{
		Id: id,
	})
	if err != nil {
		c.logger.Error("grpc DeleteSite failed", "error", err)
		return nil, err
	}
	return resp, nil
}

func (c *MonitorClient) GetSiteStatus(ctx context.Context, id string) (*monitorv1.GetSiteStatusResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.raw.GetSiteStatus(ctx, &monitorv1.GetSiteStatusRequest{
		Id: id,
	})
	if err != nil {
		c.logger.Error("grpc GetSiteStatus failed", "error", err)
		return nil, err
	}
	return resp, nil
}

func (c *MonitorClient) GetSiteHistory(
	ctx context.Context, id string, limit, offset int32,
) (*monitorv1.GetSiteHistoryResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.raw.GetSiteHistory(ctx, &monitorv1.GetSiteHistoryRequest{
		Id:     id,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		c.logger.Error("grpc GetSiteHistory failed", "error", err)
		return nil, err
	}
	return resp, nil
}
