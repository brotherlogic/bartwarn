package router_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brotherlogic/bartwarn/router"
)

const mockBARTResponse = `{
  "?xml": {
    "@version": "1.0",
    "@encoding": "utf-8"
  },
  "root": {
    "uri": {
      "#cdata-section": "http://api.bart.gov/api/sched.aspx?cmd=depart&orig=MONT&dest=PLZA&json=y"
    },
    "origin": "MONT",
    "destination": "PLZA",
    "schedule": {
      "date": "May 21, 2026",
      "time": "6:10 PM",
      "before": "2",
      "after": "2",
      "request": {
        "trip": [
          {
            "@origin": "MONT",
            "@destination": "PLZA",
            "@fare": "5.70",
            "@origTimeMin": "05:56 PM",
            "@origTimeDate": "05/21/2026",
            "@destTimeMin": "06:25 PM",
            "@destTimeDate": "05/21/2026",
            "@tripTime": "29",
            "leg": [
              {
                "@order": "1",
                "@origin": "MONT",
                "@destination": "PLZA",
                "@origTimeMin": "05:56 PM",
                "@destTimeMin": "06:25 PM",
                "@trainHeadStation": "Richmond"
              }
            ]
          },
          {
            "@origin": "MONT",
            "@destination": "PLZA",
            "@fare": "5.70",
            "@origTimeMin": "06:02 PM",
            "@origTimeDate": "05/21/2026",
            "@destTimeMin": "06:32 PM",
            "@destTimeDate": "05/21/2026",
            "@tripTime": "30",
            "leg": [
              {
                "@order": "1",
                "@origin": "MONT",
                "@destination": "19TH",
                "@origTimeMin": "06:02 PM",
                "@destTimeMin": "06:16 PM",
                "@trainHeadStation": "Antioch"
              },
              {
                "@order": "2",
                "@origin": "19TH",
                "@destination": "PLZA",
                "@origTimeMin": "06:17 PM",
                "@destTimeMin": "06:32 PM",
                "@trainHeadStation": "Richmond"
              }
            ]
          }
        ]
      }
    },
    "message": ""
  }
}`

func TestFetchTrips_Success(t *testing.T) {
	// Create a local HTTP server that returns our mock JSON
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("cmd") != "depart" || r.URL.Query().Get("orig") != "MONT" || r.URL.Query().Get("dest") != "PLZA" {
			t.Errorf("Unexpected query parameters: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockBARTResponse))
	}))
	defer ts.Close()

	// Initialize the client pointing to our test server
	client := router.NewBARTClient(ts.URL, "fake-api-key")

	// Call FetchTrips
	trips, err := client.FetchTrips(context.Background(), "MONT", "PLZA")
	if err != nil {
		t.Fatalf("FetchTrips failed: %v", err)
	}

	if len(trips) != 2 {
		t.Fatalf("Expected 2 trips, got %d", len(trips))
	}

	// Verify the first direct trip
	if trips[0].OrigTimeMin != "05:56 PM" {
		t.Errorf("Expected first trip origin time to be 05:56 PM, got %s", trips[0].OrigTimeMin)
	}
	if trips[0].DestTimeMin != "06:25 PM" {
		t.Errorf("Expected first trip dest time to be 06:25 PM, got %s", trips[0].DestTimeMin)
	}
	if len(trips[0].Legs) != 1 {
		t.Errorf("Expected first trip to have 1 leg, got %d", len(trips[0].Legs))
	}

	// Verify the second transfer trip
	if trips[1].OrigTimeMin != "06:02 PM" {
		t.Errorf("Expected second trip origin time to be 06:02 PM, got %s", trips[1].OrigTimeMin)
	}
	if len(trips[1].Legs) != 2 {
		t.Errorf("Expected second trip to have 2 legs, got %d", len(trips[1].Legs))
	}
	if trips[1].Legs[0].Destination != "19TH" {
		t.Errorf("Expected transfer at 19TH, got %s", trips[1].Legs[0].Destination)
	}
}
