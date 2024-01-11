package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aswinbennyofficial/attendease/internal/config"
	"github.com/aswinbennyofficial/attendease/internal/database"
	"github.com/aswinbennyofficial/attendease/internal/models"
	"github.com/aswinbennyofficial/attendease/internal/utility"
)

func HandleSignin(w http.ResponseWriter, r *http.Request) {

	// Instance of the Credential struct
	var creds models.Credentials
	// Get the JSON body and decode into creds
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password hash from database
	expectedPasswordHash, err := database.GetPasswordHashFromDb(creds.Username)
	if err != nil {
		log.Println("Error while getting password from database: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		if err.Error() == "User does not exist" {
			w.Write([]byte("User does not exist"))
			return
		}
		w.Write([]byte("Error while getting password from database"))
		return
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if utility.CheckPasswordHash(creds.Password, expectedPasswordHash) == false {
		log.Println("Incorrect password")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect password"))
		return
	}

	// Create a new JWT token
	signedToken, err := utility.GenerateToken(creds.Username)
	if err != nil {
		log.Println("ERROR OCCURRED WHILE CREATING JWT TOKEN: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error occurred while creating JWT token " + err.Error()))
		return
	}

	log.Printf("JWT GENERATED FOR %s", creds.Username)

	// TODO
	// Setting expiration time for cookie
	expirationTime := time.Now().Add(time.Duration(config.LoadJwtExpiresIn()) * time.Minute)

	http.SetCookie(w, &http.Cookie{
		Name:    "JWtoken",
		Value:   signedToken,
		Path:    "/",
		Expires: expirationTime,
	})

	w.Write([]byte("Login successful"))

}

func HandleRefresh(w http.ResponseWriter, r *http.Request) {
	// Parse and validate JWT from request
	claims, err := utility.ParseAndValidateJWT(r)
	if err != nil {
		log.Println("ERROR WHILE PARSING/VALIDATING JWT: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Error while parsing/validating JWT"))
		return
	}

	// Generate a new JWT token
	signedToken, err := utility.GenerateToken(claims.Username)
	if err != nil {
		log.Println("ERROR OCCURRED WHILE CREATING JWT TOKEN: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Setting JWT claims
	expirationTime := time.Now().Add(time.Duration(config.LoadJwtExpiresIn()) * time.Minute)

	http.SetCookie(w, &http.Cookie{
		Name:    "JWtoken",
		Path:    "/",
		Value:   signedToken,
		Expires: expirationTime,
	})

	log.Println("TOKEN REFRESH SUCCESSFUL")
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	log.Println("LOGOUT SUCCESSFUL")
	http.SetCookie(w, &http.Cookie{
		Name:    "JWtoken",
		Path:    "/",
		Expires: time.Now(),
	})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout successful"))
}

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	// Instance of the NewUser struct
	var user models.NewUser
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if username already exists
	isUserExist, err := database.DoesUserExist(user.Username)
	if err != nil {
		log.Println("Error while checking if user exists: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isUserExist {
		log.Println("User already exists")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("User already exists"))
		return
	}
	log.Println("User does not exist")

	// Hashing the password with the default cost of 10
	hashedPassword, err := utility.HashPassword(user.Password)
	if err != nil {
		log.Println("Error while hashing password: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while hashing password"))
		return
	}

	// Replacing existing password with hashed password
	user.Password = hashedPassword

	// Adding user and details to database
	err = database.AddUserToDb(user)
	if err != nil {
		log.Println("Error while adding user to database: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while adding user to database"))
		return
	}

	log.Println("User added to database")

	// Generate a new JWT token
	signedToken, err := utility.GenerateToken(user.Username)
	if err != nil {
		log.Println("ERROR OCCURRED WHILE CREATING JWT TOKEN: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Setting JWT claims
	expirationTime := time.Now().Add(time.Duration(config.LoadJwtExpiresIn()) * time.Minute)

	http.SetCookie(w, &http.Cookie{
		Name:    "JWtoken",
		Path:    "/",
		Value:   signedToken,
		Expires: expirationTime,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User added to database"))

}
