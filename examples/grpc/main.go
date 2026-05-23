package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/victor-mazeli-miva/go-actuator"
	grpcadapter "github.com/victor-mazeli-miva/go-actuator/adapters/grpc"
	"google.golang.org/grpc"
)

type demoCheck struct {
	name string
}

func (d demoCheck) Name() string { return d.name }

func (d demoCheck) Check(_ context.Context) error { return nil }

func main() {
	act := actuator.New()
	act.RegisterHealthCheck(demoCheck{name: "demo"})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer()
	grpcadapter.Register(srv, act)

	fmt.Println("gRPC server listening on :50051")
	fmt.Println("  grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check")
	fmt.Println("  grpcurl -plaintext -d '{\"service\":\"demo\"}' localhost:50051 grpc.health.v1.Health/Check")

	if err := srv.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
