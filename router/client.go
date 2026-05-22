package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// BARTClient interacts with the public BART API
type BARTClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewBARTClient creates a new BART API client
func NewBARTClient(baseURL, apiKey string) *BARTClient {
	if baseURL == "" {
		baseURL = "https://api.bart.gov/api/sched.aspx"
	}
	if apiKey == "" {
		apiKey = "MW9S-E7SL-26DU-VV8V" // Public demo key
	}
	return &BARTClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{},
	}
}

// Trip represents a single route option returned by BART
type Trip struct {
	OrigTimeMin string `json:"@origTimeMin"`
	DestTimeMin string `json:"@destTimeMin"`
	Legs        []Leg  `json:"leg"`
}

// Leg represents a segment of a trip (e.g. before/after a transfer)
type Leg struct {
	Destination string `json:"@destination"`
}

type bartResponse struct {
	Root struct {
		Schedule struct {
			Request struct {
				Trip []Trip `json:"trip"`
			} `json:"request"`
		} `json:"schedule"`
	} `json:"root"`
}

// FetchTrips queries the BART API for departing trips from origin to destination
func (c *BARTClient) FetchTrips(ctx context.Context, origin, destination string) ([]Trip, error) {
	reqURL := fmt.Sprintf("%s?cmd=depart&orig=%s&dest=%s&key=%s&json=y",
		c.baseURL, url.QueryEscape(origin), url.QueryEscape(destination), url.QueryEscape(c.apiKey))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected API status code: %d", resp.StatusCode)
	}

	var parsedResp bartResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsedResp); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	return parsedResp.Root.Schedule.Request.Trip, nil
}
