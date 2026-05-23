// Package grpcadapter registers the standard gRPC health checking protocol on a gRPC server.
package grpcadapter

import (
	"context"

	"github.com/victor-mazeli-miva/go-actuator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// Register registers the grpc.health.v1.Health service on srv using the actuator's health checks.
//
// Check with an empty service name reports overall health (all registered checks).
// A non-empty service name must match a HealthCheck.Name() value.
func Register(srv *grpc.Server, a *actuator.Actuator) {
	grpc_health_v1.RegisterHealthServer(srv, &healthServer{act: a})
}

type healthServer struct {
	grpc_health_v1.UnimplementedHealthServer
	act *actuator.Actuator
}

func (h *healthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	resp := h.act.EvaluateHealth(ctx)
	service := req.GetService()

	if service == "" {
		return healthResponse(resp.IsHealthy()), nil
	}

	st, ok := resp.Checks[service]
	if !ok {
		return nil, status.Error(codes.NotFound, "service not found")
	}

	return healthResponse(actuator.CheckIsHealthy(st)), nil
}

func (h *healthServer) List(ctx context.Context, _ *grpc_health_v1.HealthListRequest) (*grpc_health_v1.HealthListResponse, error) {
	resp := h.act.EvaluateHealth(ctx)
	statuses := make(map[string]*grpc_health_v1.HealthCheckResponse, len(resp.Checks)+1)
	statuses[""] = healthResponse(resp.IsHealthy())

	for name, st := range resp.Checks {
		statuses[name] = healthResponse(actuator.CheckIsHealthy(st))
	}

	return &grpc_health_v1.HealthListResponse{Statuses: statuses}, nil
}

func (h *healthServer) Watch(_ *grpc_health_v1.HealthCheckRequest, _ grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "watch not supported")
}

func healthResponse(healthy bool) *grpc_health_v1.HealthCheckResponse {
	if healthy {
		return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}
	}
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING}
}
