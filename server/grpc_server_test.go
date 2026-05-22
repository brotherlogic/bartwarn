package server_test

import (
	"context"
	"testing"

	"github.com/brotherlogic/bartwarn/api"
	"github.com/brotherlogic/bartwarn/server"
)

// MockRouter simulates the routing engine dependency
type MockRouter struct {
	called bool
}

func (m *MockRouter) FindRoute(ctx context.Context, stationId string) (string, error) {
	m.called = true
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
	router := &MockRouter{}
	notifier := &MockNotifier{}

	// Create the server injecting our mocks
	srv := server.NewBartwarnServer(router, notifier)

	req := &api.LocationRequest{
		StationId: "MONT",
	}

	// Call the gRPC method
	_, err := srv.RecordLocation(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected error from RecordLocation: %v", err)
	}

	// Verify our dependencies were invoked
	if !router.called {
		t.Errorf("Expected router.FindRoute to be called")
	}
	if !notifier.called {
		t.Errorf("Expected notifier.SendSMS to be called")
	}
}
