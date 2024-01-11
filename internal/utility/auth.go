package utility

import (
	"errors"

	"net/http"

	"time"

	"github.com/aswinbennyofficial/attendease/internal/config"
	"github.com/aswinbennyofficial/attendease/internal/models"
	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// Hashing the password with the default cost of 10
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Parse and validate JWT from request
func ParseAndValidateJWT(r *http.Request) (*models.Claims, error) {
	cookieVar, err := r.Cookie("JWtoken")
	if err != nil {
		return nil, err
	}

	// Get JWT string from cookie
	JWTstring := cookieVar.Value

	// Initialize a new instance of Claims
	claims := &models.Claims{}

	// Parse the JWT string and store in claims
	tokenVar, err := jwt.ParseWithClaims(JWTstring, claims, func(token *jwt.Token) (interface{}, error) {
		return config.LoadJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if !tokenVar.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// Generate a new JWT token
func GenerateToken(username string) (string, error) {
	// Setting expiration time to be 5 minutes from now
	expirationTime := time.Now().Add(time.Duration(config.LoadJwtExpiresIn()) * time.Minute)

	claims := &models.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declaring token with header and payload
	noSignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create a complete signed JWT
	signedToken, err := noSignedToken.SignedString(config.LoadJWTSecret())
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
