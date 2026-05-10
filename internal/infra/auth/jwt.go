package auth

import (
	"os"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

var TokenAuth *jwtauth.JWTAuth

func GenerateJWT() {
	_ = godotenv.Load()
	secret := os.Getenv("JWT_SECRET")

	TokenAuth = jwtauth.New(
		"HS256",
		[]byte(secret),
		nil,
		jwt.WithAcceptableSkew(30*time.Second),
	)
}
