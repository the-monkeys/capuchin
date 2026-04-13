package middleware_test

// Tests for DBHealthCheck middleware.
//
// Feature: migration-module-separation, Property 8: 503 returned while unhealthy.
// No Docker or real Postgres required — health state is set directly via the
// database package's exported atomic flag.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"capuchin/internal/database"
	"capuchin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// setDBHealthy directly manipulates the package-level atomic flag for testing.
// We access it via the exported IsHealthy() path by reflecting the flag through
// the database package — but since dbHealthy is unexported, we drive it via
// the exported pointer trick: expose a test helper in the database package.
// Instead, we use the simpler approach: call database.ForceHealthyForTest.
// Since that doesn't exist, we drive the flag indirectly by calling
// atomic.StoreInt32 on the exported symbol via a test-only shim.
//
// The cleanest approach without modifying production code: use a build-tag
// test helper. Here we use the internal test package approach — the middleware
// test is in package middleware_test (external), so we drive health state via
// the database package's exported SetHealthForTest helper if available, or
// we accept that the flag starts at 0 (unhealthy) and test both states by
// directly importing the atomic.

// forceHealth sets the database health flag for testing purposes.
// It accesses the unexported dbHealthy via the internal test package trick:
// since this file is in middleware_test (external package), we use a small
// exported shim in the database package's test file.
// For simplicity, we use the database_test-accessible atomic directly.
func forceHealth(healthy bool) {
	if healthy {
		atomic.StoreInt32(database.DBHealthyPtr(), 1)
	} else {
		atomic.StoreInt32(database.DBHealthyPtr(), 0)
	}
}

func TestDBHealthMiddleware_Returns503WhenUnhealthy(t *testing.T) {
	forceHealth(false)
	t.Cleanup(func() { forceHealth(false) })

	r := gin.New()
	r.Use(middleware.DBHealthCheck())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}

	var body map[string]string
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	if body["error"] != "database unavailable" {
		t.Errorf("expected error='database unavailable', got %q", body["error"])
	}
}

func TestDBHealthMiddleware_PassesWhenHealthy(t *testing.T) {
	forceHealth(true)
	t.Cleanup(func() { forceHealth(false) })

	handlerCalled := false
	r := gin.New()
	r.Use(middleware.DBHealthCheck())
	r.GET("/test", func(c *gin.Context) {
		handlerCalled = true
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !handlerCalled {
		t.Error("expected handler to be called when DB is healthy")
	}
}
