package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name           string
		input          User
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Credentials",
			input:          User{Username: "test", Password: "password"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"token": "fake-jwt-token"}`,
		},
		{
			name:           "Invalid Username",
			input:          User{Username: "wrong", Password: "password"},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid credentials\n",
		},
		{
			name:           "Invalid Password",
			input:          User{Username: "test", Password: "wrong"},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid credentials\n",
		},
		{
			name:           "Invalid JSON",
			input:          User{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid input\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if tt.name != "Invalid JSON" {
				var err error
				body, err = json.Marshal(tt.input)
				if err != nil {
					t.Fatalf("Failed to marshal input: %v", err)
				}
			} else {
				body = []byte("invalid-json")
			}

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			LoginHandler(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			responseBody := w.Body.String()
			if responseBody != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, responseBody)
			}
		})
	}
}
