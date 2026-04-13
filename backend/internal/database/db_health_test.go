package database

// TestDBHealthFlag_ReflectsConnectionState verifies that IsHealthy() reflects
// the atomic dbHealthy flag correctly without requiring a real Postgres instance.
//
// Feature: migration-module-separation, Property 7: DB health flag reflects connection state.

import (
	"sync/atomic"
	"testing"
)

func TestDBHealthFlag_ReflectsConnectionState(t *testing.T) {
	// Reset state after test.
	t.Cleanup(func() { atomic.StoreInt32(&dbHealthy, 0) })

	// Initially unhealthy.
	atomic.StoreInt32(&dbHealthy, 0)
	if IsHealthy() {
		t.Error("expected IsHealthy() == false when dbHealthy=0")
	}

	// Simulate successful ping — mark healthy.
	atomic.StoreInt32(&dbHealthy, 1)
	if !IsHealthy() {
		t.Error("expected IsHealthy() == true when dbHealthy=1")
	}

	// Simulate lost connection — mark unhealthy.
	atomic.StoreInt32(&dbHealthy, 0)
	if IsHealthy() {
		t.Error("expected IsHealthy() == false after dbHealthy reset to 0")
	}

	// Simulate recovery.
	atomic.StoreInt32(&dbHealthy, 1)
	if !IsHealthy() {
		t.Error("expected IsHealthy() == true after recovery")
	}
}
