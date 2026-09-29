package client

import (
	"fmt"

	monitorv1 "gitlab.com/Dokuchaevvn/site-monitor/gen/go/monitor/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewMonitorClient(addr string) (monitorv1.MonitorServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("dial monitor grpc %s: %w", addr, err)
	}

	return monitorv1.NewMonitorServiceClient(conn), conn, nil
}
