package routes_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"discount-system-backend/internal/routes"
)

func TestRoutes(t *testing.T) {
	mux := routes.Routes()

	t.Run("ping endpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/ping", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}

		expectedBody := "pong"
		if strings.TrimSpace(rr.Body.String()) != expectedBody {
			t.Errorf("expected body %q, got %q", expectedBody, rr.Body.String())
		}
	})

	t.Run("login endpoint", func(t *testing.T) {
		reqBody := `{"username":"test","password":"password"}`
		req := httptest.NewRequest("POST", "/login", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}
