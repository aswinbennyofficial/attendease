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

func HandleAdminSignin(w http.ResponseWriter, r *http.Request) {

	// Instance of the Credential struct
	var creds models.LoginCreds
	// Get the JSON body and decode into creds
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}


	// Get the expected password hash from database
	expectedUsername ,expectedPasswordHash, err := database.GetHashAndUsernameFromDb(creds.Organisation)
	if err != nil {
		log.Println("Error while getting password from database: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		if err.Error() == "Organisation does not exist" {
			w.Write([]byte("Organisation does not exist"))
			return
		}
		w.Write([]byte("Error while getting password from database"))
		return
	}

	if expectedUsername != creds.Username {
		log.Println("Incorrect username")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect username"))
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
	signedToken, err := utility.GenerateToken(creds.Organisation,creds.Username,true,false)
	if err != nil {
		log.Println("ERROR OCCURRED WHILE CREATING JWT TOKEN: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error occurred while creating JWT token " + err.Error()))
		return
	}

	log.Printf("JWT GENERATED FOR %s", creds.Organisation)

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
	signedToken, err := utility.GenerateToken(claims.Org,claims.Username,claims.Admin,claims.Employee)
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

func HandleAdminSignup(w http.ResponseWriter, r *http.Request) {
	// Instance of the NewUser struct
	var org models.NewUser
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&org)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if org.Organisation==""{
		log.Println("Organisation cannot be empty")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Organisation cannot be empty"))
		return
	}
	
	if org.Username==""{
		log.Println("Username cannot be empty")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Username cannot be empty"))
		return
	}

	// Check if organisation already exists
	isOrgExist, err := database.DoesExistInAuthColl("organisation",org.Organisation)
	if err != nil {
		log.Println("Error while checking if user exists: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isOrgExist {
		log.Println("Organisation already exists")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Organisation already exists"))
		return
	}
	log.Println("Organisation does not exist")

	// Check if username already exists
	isUsernameExist, err := database.DoesExistInAuthColl("username",org.Username)
	if err != nil {
		log.Println("Error while checking if user exists: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if isUsernameExist {
		log.Println("Username already exists")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("Username already exists"))
		return
	}

	// Hashing the password with the default cost of 10
	hashedPassword, err := utility.HashPassword(org.Password)
	if err != nil {
		log.Println("Error while hashing password: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while hashing password"))
		return
	}

	// Replacing existing password with hashed password
	org.Password = hashedPassword

	org.IsVerified = false
	org.VerifyCode,err = utility.CreateRandomString(15)

	if err != nil {	
		log.Println("Error while creating random string: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while creating verification string"))
		return
	}

	// Adding user and details to database
	err = database.AddUserToDb(org)
	if err != nil {
		log.Println("Error while adding Oreganisation to database: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while adding Organisation to database"))
		return
	}

	log.Println("Organisation added to database")

	// Generate a new JWT token
	signedToken, err := utility.GenerateToken(org.Organisation,org.Username,true,false)
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
	w.Write([]byte("Organisation signup successful"))

}
