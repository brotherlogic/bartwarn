package router

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

// ErrSuppressSMS indicates that the system should quietly ignore this ping
var ErrSuppressSMS = errors.New("suppress sms due to cache ttl or rules")

// TripFetcher abstracts the API client for testing
type TripFetcher interface {
	FetchTrips(ctx context.Context, origin, destination string) ([]Trip, error)
}

// Engine implements the core business logic and caching
type Engine struct {
	fetcher     TripFetcher
	destination string
	mu          sync.Mutex
	cacheTTL    time.Time
}

// NewEngine creates a new routing engine
func NewEngine(fetcher TripFetcher, destination string) *Engine {
	return &Engine{
		fetcher:     fetcher,
		destination: destination,
	}
}

// FindRoute evaluates the trips and returns the SMS string if a warning is needed
func (e *Engine) FindRoute(ctx context.Context, stationId string) (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	now := time.Now()

	// 1. Check Global Cache TTL
	if now.Before(e.cacheTTL) {
		slog.Debug("Suppressing ping due to TTL cache", "until", e.cacheTTL)
		return "", ErrSuppressSMS
	}

	trips, err := e.fetcher.FetchTrips(ctx, stationId, e.destination)
	if err != nil {
		return "", fmt.Errorf("failed to fetch trips: %w", err)
	}

	// 2. Find best valid trip
	for _, trip := range trips {
		origTime, err := time.Parse("03:04 PM", trip.OrigTimeMin)
		if err != nil {
			continue
		}

		// Normalize to today's date for accurate comparison
		todayOrigTime := time.Date(now.Year(), now.Month(), now.Day(), origTime.Hour(), origTime.Minute(), 0, 0, now.Location())

		// Handle midnight wrapping
		if todayOrigTime.Before(now) && now.Sub(todayOrigTime) > 12*time.Hour {
			todayOrigTime = todayOrigTime.Add(24 * time.Hour)
		}

		// 3. Apply the 2-minute grace period filter
		if todayOrigTime.Sub(now) >= 2*time.Minute {
			// Valid trip found! Set cache TTL to its arrival time.
			destTime, err := time.Parse("03:04 PM", trip.DestTimeMin)
			if err == nil {
				todayDestTime := time.Date(now.Year(), now.Month(), now.Day(), destTime.Hour(), destTime.Minute(), 0, 0, now.Location())
				if todayDestTime.Before(todayOrigTime) {
					todayDestTime = todayDestTime.Add(24 * time.Hour)
				}
				e.cacheTTL = todayDestTime
			} else {
				// Fallback TTL
				e.cacheTTL = now.Add(30 * time.Minute)
			}

			// 4. Format SMS Message
			var transferInfo string
			if len(trip.Legs) > 1 {
				transferInfo = fmt.Sprintf(" (transfer at %s)", trip.Legs[0].Destination)
			}

			msg := fmt.Sprintf("Next trip to %s departs at %s%s. Arrives at %s.",
				e.destination, trip.OrigTimeMin, transferInfo, trip.DestTimeMin)

			return msg, nil
		}
	}

	// No valid trains found
	return "", ErrSuppressSMS
}
