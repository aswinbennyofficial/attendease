package config

import (
	"os"
	"log"
	"strconv"
)

func LoadJWTSecret() []byte{
	JWT_SECRET:= os.Getenv("JWT_SECRET")
	if JWT_SECRET==""{
		log.Println("Error loading JWT_SECRET in config.LoadJWTSecret()")
		return []byte("")
	}
	return []byte(JWT_SECRET)
}

func LoadJwtExpiresIn() int {

	jwtExpiresInStr := os.Getenv("JWT_EXPIRES_IN")
	if jwtExpiresInStr == "" {
		log.Println("JWT_EXPIRES_IN is not set, using default value.")
		return 5 // default value 
	}

	jwtExpiresIn, err := strconv.Atoi(jwtExpiresInStr)
	if err != nil {
		log.Printf("Error converting JWT_EXPIRES_IN to integer: %v\n", err)
		return 5 // default value 
	}

	return jwtExpiresIn
}

