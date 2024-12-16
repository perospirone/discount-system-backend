package jwt

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type CustomClaims struct {
	Name     string `json:"name"`
	Password string `json:"email"`
	jwt.StandardClaims
}

var Secret = []byte(os.Getenv("DB_HOST")) // TODO: make .env

func CreateTokenJWT(name, email string) (string, error) {
	claims := &CustomClaims{
		name,
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(Secret)
	if err != nil {
		log.Println("Error: ", err)
	}

	return t, err
}
