package auth_test

import (
	"bytes"
	"discount-system-backend/internal/auth"
	"discount-system-backend/internal/database"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginSuccessToken(t *testing.T) {
	reqBody, _ := json.Marshal(map[string]string{
		"name": "test",
		"password": "password",
	})

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.LoginHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status code %d but got %d", http.StatusOK, status)
	}

	var response auth.ResponseToken

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response body into User struct: %v", err)
	}

	if response.Token == "" {
		t.Errorf("expected response token but got empty value")
	}
}

func TestLoginHandler(t *testing.T) { // change and divide this test in more tests
	tests := []struct {
		name           string
		input          database.User
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid Credentials",
			input:          database.User{Name: "test", Password: "password"},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"token": "fake-jwt-token"}`,
		},
		{
			name:           "Invalid Username",
			input:          database.User{Name: "wrong", Password: "password"},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid credentials\n",
		},
		{
			name:           "Invalid Password",
			input:          database.User{Name: "test", Password: "wrong"},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid credentials\n",
		},
		{
			name:           "Invalid JSON",
			input:          database.User{},
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

			auth.LoginHandler(w, req)

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
