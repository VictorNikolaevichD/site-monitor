package client

import (
	"context"
	"fmt"
	"time"

	monitorv1 "gitlab.com/Dokuchaevvn/site-monitor/gen/go/monitor/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

func NewMonitorClient(addr string) (monitorv1.MonitorServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
			return nil, nil, fmt.Errorf("monitor grpc %s not ready %w", addr, ctx.Err())
		}
	}

	return monitorv1.NewMonitorServiceClient(conn), conn, nil
}
