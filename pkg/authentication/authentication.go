package authentication

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(secret string, claims Claims) (token string, err error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	return jwtToken.SignedString([]byte(secret))
}
