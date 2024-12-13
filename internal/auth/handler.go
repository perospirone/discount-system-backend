package auth

import (
	"encoding/json"
	"net/http"
)

// Example user struct (replace with DB integration)
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseToken struct {
	Token string `json:"token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate user (example, replace with real validation)
	if user.Username == "test" && user.Password == "password" {
		// Return JWT (stub response)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"token": "fake-jwt-token"}`))
		//json.NewEncoder(w).Encode(ResponseToken{Token: "fake-jwt-token"})
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}
