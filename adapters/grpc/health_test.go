package grpcadapter

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/victor-mazeli-miva/go-actuator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

type stubCheck struct {
	name string
	err  error
}

func (s stubCheck) Name() string { return s.name }

func (s stubCheck) Check(_ context.Context) error { return s.err }

func startServer(t *testing.T, act *actuator.Actuator) (*grpc.ClientConn, func()) {
	t.Helper()

	lis := bufconn.Listen(bufSize)
	srv := grpc.NewServer()
	Register(srv, act)

	go func() {
		if err := srv.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			t.Errorf("serve: %v", err)
		}
	}()

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}

	cleanup := func() {
		_ = conn.Close()
		srv.Stop()
	}

	return conn, cleanup
}

func TestCheckOverallServing(t *testing.T) {
	act := actuator.New()
	act.RegisterHealthCheck(stubCheck{name: "postgres"})
	act.RegisterHealthCheck(stubCheck{name: "redis"})

	conn, cleanup := startServer(t, act)
	defer cleanup()

	client := grpc_health_v1.NewHealthClient(conn)
	resp, err := client.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Fatalf("status = %v, want SERVING", resp.Status)
	}
}

func TestCheckOverallNotServing(t *testing.T) {
	act := actuator.New()
	act.RegisterHealthCheck(stubCheck{name: "postgres", err: errors.New("down")})
	act.RegisterHealthCheck(stubCheck{name: "redis"})

	conn, cleanup := startServer(t, act)
	defer cleanup()

	client := grpc_health_v1.NewHealthClient(conn)
	resp, err := client.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if resp.Status != grpc_health_v1.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("status = %v, want NOT_SERVING", resp.Status)
	}
}

func TestCheckNamedService(t *testing.T) {
	act := actuator.New()
	act.RegisterHealthCheck(stubCheck{name: "postgres", err: errors.New("down")})
	act.RegisterHealthCheck(stubCheck{name: "redis"})

	conn, cleanup := startServer(t, act)
	defer cleanup()

	client := grpc_health_v1.NewHealthClient(conn)

	resp, err := client.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "redis"})
	if err != nil {
		t.Fatalf("Check redis: %v", err)
	}
	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Fatalf("redis status = %v, want SERVING", resp.Status)
	}

	resp, err = client.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "postgres"})
	if err != nil {
		t.Fatalf("Check postgres: %v", err)
	}
	if resp.Status != grpc_health_v1.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("postgres status = %v, want NOT_SERVING", resp.Status)
	}
}

func TestCheckUnknownService(t *testing.T) {
	act := actuator.New()

	conn, cleanup := startServer(t, act)
	defer cleanup()

	client := grpc_health_v1.NewHealthClient(conn)
	_, err := client.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "missing"})
	if err == nil {
		t.Fatal("expected error for unknown service")
	}
	if status.Code(err) != codes.NotFound {
		t.Fatalf("code = %v, want NotFound", status.Code(err))
	}
}

func TestList(t *testing.T) {
	act := actuator.New()
	act.RegisterHealthCheck(stubCheck{name: "postgres"})
	act.RegisterHealthCheck(stubCheck{name: "redis", err: errors.New("down")})

	conn, cleanup := startServer(t, act)
	defer cleanup()

	client := grpc_health_v1.NewHealthClient(conn)
	resp, err := client.List(context.Background(), &grpc_health_v1.HealthListRequest{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	if resp.Statuses[""].Status != grpc_health_v1.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("overall = %v, want NOT_SERVING", resp.Statuses[""].Status)
	}
	if resp.Statuses["postgres"].Status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Fatalf("postgres = %v, want SERVING", resp.Statuses["postgres"].Status)
	}
	if resp.Statuses["redis"].Status != grpc_health_v1.HealthCheckResponse_NOT_SERVING {
		t.Fatalf("redis = %v, want NOT_SERVING", resp.Statuses["redis"].Status)
	}
}

func TestWatchUnimplemented(t *testing.T) {
	act := actuator.New()

	conn, cleanup := startServer(t, act)
	defer cleanup()

	client := grpc_health_v1.NewHealthClient(conn)
	stream, err := client.Watch(context.Background(), &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		if status.Code(err) == codes.Unimplemented {
			return
		}
		t.Fatalf("Watch: %v", err)
	}

	_, err = stream.Recv()
	if err == nil {
		t.Fatal("expected error from Watch stream")
	}
	if status.Code(err) != codes.Unimplemented {
		t.Fatalf("code = %v, want Unimplemented", status.Code(err))
	}
}
