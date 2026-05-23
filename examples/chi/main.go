package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/victor-mazeli-miva/go-actuator"
	"github.com/victor-mazeli-miva/go-actuator/adapters"
)

type demoCheck struct {
	name string
}

func (d demoCheck) Name() string { return d.name }

func (d demoCheck) Check(_ context.Context) error { return nil }

func main() {
	act := actuator.New()
	act.RegisterHealthCheck(demoCheck{name: "demo"})

	r := chi.NewRouter()
	adapters.MountChi(r, "/actuator", act)

	fmt.Println("Actuator listening on :8080")
	fmt.Println("  curl http://localhost:8080/actuator/health")
	fmt.Println("  curl http://localhost:8080/actuator/live")
	fmt.Println("  curl http://localhost:8080/actuator/ready")
	fmt.Println("  curl http://localhost:8080/actuator/runtime")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
