package router_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brotherlogic/bartwarn/router"
)

type MockClient struct {
	trips []router.Trip
}

func (m *MockClient) FetchTrips(ctx context.Context, origin, destination string) ([]router.Trip, error) {
	return m.trips, nil
}

func TestFindRoute_SuccessAndCache(t *testing.T) {
	now := time.Now()
	
	// Trip departing in 1 min (should be skipped due to 2 min grace period)
	trip1 := router.Trip{
		OrigTimeMin: now.Add(1 * time.Minute).Format("03:04 PM"),
		DestTimeMin: now.Add(30 * time.Minute).Format("03:04 PM"),
		Legs:        []router.Leg{{Destination: "PLZA"}},
	}
	
	// Trip departing in 5 mins (should be selected)
	trip2 := router.Trip{
		OrigTimeMin: now.Add(5 * time.Minute).Format("03:04 PM"),
		DestTimeMin: now.Add(35 * time.Minute).Format("03:04 PM"),
		Legs:        []router.Leg{{Destination: "19TH"}, {Destination: "PLZA"}},
	}

	mockClient := &MockClient{trips: []router.Trip{trip1, trip2}}
	
	engine := router.NewEngine(mockClient, "PLZA")

	// Call FindRoute (First ping)
	msg, err := engine.FindRoute(context.Background(), "MONT")
	if err != nil {
		t.Fatalf("Unexpected error on first ping: %v", err)
	}

	expectedPrefix := "Next trip to PLZA departs at " + trip2.OrigTimeMin
	if msg == "" || msg[:len(expectedPrefix)] != expectedPrefix {
		t.Errorf("Expected message to start with %q, got %q", expectedPrefix, msg)
	}

	// Call FindRoute (Second ping immediately after)
	_, err = engine.FindRoute(context.Background(), "MONT")
	if !errors.Is(err, router.ErrSuppressSMS) {
		t.Errorf("Expected ErrSuppressSMS due to cache TTL, got %v", err)
	}
}
