package auth

import (
	"discount-system-backend/internal/database"
	"discount-system-backend/pkg/jwt"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

type ResponseToken struct {
	Token string `json:"token"`
}

type responseTokenAndID struct {
	Token string `json:"token"`
	UserID uint `json:"user_id"`
}

type responseError struct {
	Error string `json:"error"`
}

var db = database.Connection()

//func LoginHandler(w http.ResponseWriter, r *http.Request) {
//var user database.User
//if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
//http.Error(w, "Invalid input", http.StatusBadRequest)
//return
//}

//// Validate user (example, replace with real validation)
//if user.Name == "test" && user.Password == "password" {
//// Return JWT (stub response)
//w.WriteHeader(http.StatusOK)
//w.Write([]byte(`{"token": "fake-jwt-token"}`))
////json.NewEncoder(w).Encode(ResponseToken{Token: "fake-jwt-token"})
//} else {
//http.Error(w, "Invalid credentials", http.StatusUnauthorized)
//}
//}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	body := &database.User{}
	user := &database.User{}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read request body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseError{Error: "Internal server error"})
		return
	}

	err = json.Unmarshal(reqBody, &body)
	if err != nil {
		log.Println("Failed to unmarshal JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseError{Error: "Invalid JSON format"})
		return
	}

	result := db.Take(&user, "email = ?", body.Email)
	if result.Error != nil {
		log.Println("User not found or database error:", result.Error)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(responseError{Error: "Invalid email or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(responseError{Error: "Invalid email or password"})
		return
	}

	token, err := jwt.CreateTokenJWT(user.Name, user.Email)
	if err != nil {
		log.Println("Failed to create JWT:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseError{Error: "Failed to generate token"})
		return
	}

	response := responseTokenAndID{Token: token, UserID: user.ID}
	json.NewEncoder(w).Encode(response)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	body := &database.User{}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read request body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseError{Error: "Internal server error"})
		return
	}

	err = json.Unmarshal(reqBody, &body)
	if err != nil {
		log.Println("Failed to unmarshal JSON:", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseError{Error: "Invalid JSON format"})
		return
	}

	if _, err := mail.ParseAddress(body.Email); err != nil {
		log.Println("Failed to parse email", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseError{Error: "Invalid Email"})
		return
	}

	if body.Name == "" {
		log.Println("empty name")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseError{Error: "Invalid Name"})
		return
	}

	// Check if the user already exists by email
	var existingUser database.User
	result := db.Take(&existingUser, "email = ?", body.Email)
	if result.Error == nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(responseError{Error: "User already exists"})
		return
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to hash password:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseError{Error: "Failed to process password"})
		return
	}

	// Store the hashed password instead of plain text
	body.Password = string(hashedPassword)

	// Save the new user to the database
	result = db.Create(&body)
	if result.Error != nil {
		log.Println("Failed to create user in database:", result.Error)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseError{Error: "Failed to register user"})
		return
	}

	token, err := jwt.CreateTokenJWT(body.Name, body.Email)
	if err != nil {
		log.Println("Failed to create JWT:", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseError{Error: "Failed to generate token"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := responseTokenAndID{Token: token, UserID: body.ID}
	json.NewEncoder(w).Encode(response)
}
