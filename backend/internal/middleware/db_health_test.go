package middleware_test

// Tests for DBHealthCheck middleware.
//
// Feature: migration-module-separation, Property 8: 503 returned while unhealthy.
// No Docker or real Postgres required — health state is set directly via
// database.SetHealthForTest.

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"capuchin/internal/database"
	"capuchin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestDBHealthMiddleware_Returns503WhenUnhealthy(t *testing.T) {
	database.SetHealthForTest(false)
	t.Cleanup(func() { database.SetHealthForTest(false) })

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
	database.SetHealthForTest(true)
	t.Cleanup(func() { database.SetHealthForTest(false) })

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
