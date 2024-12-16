package routes

import (
	"discount-system-backend/internal/auth"

	"net/http"
)

func ping(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`pong`))
}

func Routes() *http.ServeMux {
 mux := http.NewServeMux()

	mux.HandleFunc("POST /login", auth.LoginHandler)
	mux.HandleFunc("/ping", ping)

	return mux
}
