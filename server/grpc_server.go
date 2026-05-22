package server

import (
	"context"
	"errors"
	"log/slog"

	"github.com/brotherlogic/bartwarn/api"
	"github.com/brotherlogic/bartwarn/router"
)

// Router represents the routing engine module (Issue #4)
type Router interface {
	FindRoute(ctx context.Context, stationId string) (string, error)
}

// Notifier represents the SMS notification module (Issue #5)
type Notifier interface {
	SendSMS(message string) error
}

// BartwarnServer implements the api.LocationServiceServer gRPC interface
type BartwarnServer struct {
	api.UnimplementedLocationServiceServer
	router   Router
	notifier Notifier
}

// NewBartwarnServer creates a new instance of the gRPC server with injected dependencies
func NewBartwarnServer(router Router, notifier Notifier) *BartwarnServer {
	return &BartwarnServer{
		router:   router,
		notifier: notifier,
	}
}

// RecordLocation handles incoming pings from the Android client
func (s *BartwarnServer) RecordLocation(ctx context.Context, req *api.LocationRequest) (*api.LocationResponse, error) {
	slog.Info("Received RecordLocation ping", 
		"station_id", req.StationId, 
		"lat", req.Latitude, 
		"long", req.Longitude,
	)

	// Delegate to the routing engine
	msg, err := s.router.FindRoute(ctx, req.StationId)
	if err != nil {
		if errors.Is(err, router.ErrSuppressSMS) {
			slog.Info("Suppressing SMS due to cache TTL or grace period rules")
			return &api.LocationResponse{}, nil
		}
		slog.Error("Router failed to find route", "error", err)
		return nil, err
	}

	// Delegate to the SMS notification client
	err = s.notifier.SendSMS(msg)
	if err != nil {
		slog.Error("Notifier failed to send SMS", "error", err)
		return nil, err
	}

	return &api.LocationResponse{}, nil
}
