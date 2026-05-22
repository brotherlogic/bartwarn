package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/brotherlogic/bartwarn/api"
	"github.com/brotherlogic/bartwarn/server"
	"google.golang.org/grpc"
)

// Dummy implementations for unbuilt modules
type dummyRouter struct{}
func (r *dummyRouter) FindRoute(stationId string) error {
	slog.Info("Dummy Router called", "stationId", stationId)
	return nil
}

type dummyNotifier struct{}
func (n *dummyNotifier) SendSMS(message string) error {
	slog.Info("Dummy Notifier called", "message", message)
	return nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		slog.Error("Failed to listen on TCP port", "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	
	router := &dummyRouter{}
	notifier := &dummyNotifier{}
	
	bartwarnServer := server.NewBartwarnServer(router, notifier)
	api.RegisterLocationServiceServer(grpcServer, bartwarnServer)

	slog.Info("Starting gRPC Server", "port", port)
	if err := grpcServer.Serve(listener); err != nil {
		slog.Error("Failed to serve gRPC", "error", err)
		os.Exit(1)
	}
}
