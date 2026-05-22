package api_test

import (
	"testing"

	"github.com/brotherlogic/bartwarn/api"
)

func TestLocationRequestInitialization(t *testing.T) {
	// Our test confirms that the protobuf-generated struct exists and can be initialized
	req := &api.LocationRequest{
		StationId: "MONT",
		Latitude:  37.789256,
		Longitude: -122.401407,
	}

	if req.StationId != "MONT" {
		t.Errorf("Expected station ID to be MONT, got %s", req.StationId)
	}
}
