package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/uLesson-Education/go-actuator"
	"github.com/uLesson-Education/go-actuator/adapters"
)

type demoCheck struct {
	name string
}

func (d demoCheck) Name() string { return d.name }

func (d demoCheck) Check(_ context.Context) error { return nil }

func main() {
	gin.SetMode(gin.ReleaseMode)

	act := actuator.New()
	act.RegisterHealthCheck(demoCheck{name: "demo"})

	r := gin.New()
	adapters.MountGin(r, "/actuator", act)

	fmt.Println("Actuator listening on :8080")
	fmt.Println("  curl http://localhost:8080/actuator/health")
	fmt.Println("  curl http://localhost:8080/actuator/live")
	fmt.Println("  curl http://localhost:8080/actuator/ready")
	fmt.Println("  curl http://localhost:8080/actuator/runtime")

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
