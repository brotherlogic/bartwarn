package server_test

import (
	"context"
	"testing"

	"github.com/brotherlogic/bartwarn/api"
	"github.com/brotherlogic/bartwarn/router"
	"github.com/brotherlogic/bartwarn/server"
)

// MockRouter simulates the routing engine dependency
type MockRouter struct {
	called       bool
	suppressPing bool
}

func (m *MockRouter) FindRoute(ctx context.Context, stationId string) (string, error) {
	m.called = true
	if m.suppressPing {
		return "", router.ErrSuppressSMS
	}
	return "Mock SMS Message", nil
}

// MockNotifier simulates the SMS notification dependency
type MockNotifier struct {
	called bool
}

func (m *MockNotifier) SendSMS(message string) error {
	m.called = true
	return nil
}

func TestRecordLocation_Success(t *testing.T) {
	r := &MockRouter{}
	n := &MockNotifier{}

	srv := server.NewBartwarnServer(r, n)
	req := &api.LocationRequest{StationId: "MONT"}

	_, err := srv.RecordLocation(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error from RecordLocation: %v", err)
	}

	if !r.called {
		t.Errorf("Expected router to be called")
	}
	if !n.called {
		t.Errorf("Expected notifier to be called")
	}
}

func TestRecordLocation_SuppressSMS(t *testing.T) {
	r := &MockRouter{suppressPing: true}
	n := &MockNotifier{}

	srv := server.NewBartwarnServer(r, n)
	req := &api.LocationRequest{StationId: "MONT"}

	_, err := srv.RecordLocation(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error from RecordLocation when suppressing: %v", err)
	}

	if !r.called {
		t.Errorf("Expected router to be called")
	}
	if n.called {
		t.Errorf("Expected notifier NOT to be called")
	}
}
