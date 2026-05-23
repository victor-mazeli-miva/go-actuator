package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/victor-mazeli-miva/go-actuator"
)

type demoCheck struct {
	name string
}

func (d demoCheck) Name() string { return d.name }

func (d demoCheck) Check(_ context.Context) error { return nil }

func main() {
	act := actuator.New()
	act.RegisterHealthCheck(demoCheck{name: "demo"})

	fmt.Println("Actuator listening on :8080")
	fmt.Println("  curl http://localhost:8080/health")
	fmt.Println("  curl http://localhost:8080/live")
	fmt.Println("  curl http://localhost:8080/ready")
	fmt.Println("  curl http://localhost:8080/runtime")

	if err := http.ListenAndServe(":8080", act.Router()); err != nil {
		log.Fatal(err)
	}
}
