package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/brotherlogic/bartwarn/api"
	"github.com/brotherlogic/bartwarn/notifier"
	"github.com/brotherlogic/bartwarn/server"
	"google.golang.org/grpc"
)

// Dummy implementations for unbuilt modules
type dummyRouter struct{}
func (r *dummyRouter) FindRoute(ctx context.Context, stationId string) (string, error) {
	slog.Info("Dummy Router called", "stationId", stationId)
	return "Dummy message from main.go", nil
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

	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		slog.Warn("SMTP_HOST is not set; SMS notification will likely fail")
	}

	notifierClient := notifier.NewSMTPNotifier(
		smtpHost,
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
		os.Getenv("TARGET_SMS_EMAIL"),
		os.Getenv("SMTP_FROM_EMAIL"),
	)
	
	bartwarnServer := server.NewBartwarnServer(router, notifierClient)
	api.RegisterLocationServiceServer(grpcServer, bartwarnServer)

	// Start server in a goroutine so it doesn't block
	go func() {
		slog.Info("Starting gRPC Server", "port", port)
		if err := grpcServer.Serve(listener); err != nil {
			slog.Error("Failed to serve gRPC", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for OS interrupt signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down gRPC server gracefully...")
	grpcServer.GracefulStop()
	slog.Info("Server stopped")
}
