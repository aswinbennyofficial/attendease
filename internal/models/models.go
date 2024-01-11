package models

import(
	"github.com/golang-jwt/jwt/v5"
)



// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Struct to store the user data in database (signup)
type NewUser struct{
	Username string `json:"username"`
	Password string `json:"password"`
} 