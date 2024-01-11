package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aswinbennyofficial/attendease/internal/models"

)


func HandleWelcome(w http.ResponseWriter, r *http.Request){
	// Get claims from context
	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		log.Println("Claims not found in context")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Claims not found in context"))
		return
	}

	log.Println("ACCESS APPROVED TO /WELCOME")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))

		
}

func HandleHealth(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("hello worlds"))
}