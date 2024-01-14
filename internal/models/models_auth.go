package models

import(
	"github.com/golang-jwt/jwt/v5"
)



// Create a struct to read the username and password from the request body
type LoginCreds struct {
	Password string `json:"password"`
	Username string `json:"username"`
	Organisation string `json:"organisation"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	Org string `json:"org"`
	Admin bool `json:"admin"`
	Employee bool `json:"employee"`
	jwt.RegisteredClaims
}

// Struct to store the user data in database (signup)
type NewUser struct{
	Organisation string `json:"organisation"`
	Username string `json:"username"`
	Password string `json:"password"`
	
	// For email verification : optional, 
	IsVerified bool `json:"isverified"`
	VerifyCode string `json:"verifycode"`
}





