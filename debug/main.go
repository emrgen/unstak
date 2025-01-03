package main

import (
	"github.com/emrgen/unpost/internal/server"
	"os"
)

func main() {
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "4020"
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "4021"
	}

	err := server.Start(grpcPort, httpPort)
	if err != nil {
		return
	}
}
