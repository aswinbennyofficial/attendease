package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aswinbennyofficial/attendease/internal/config"
	"github.com/aswinbennyofficial/attendease/internal/database"
	
	"github.com/aswinbennyofficial/attendease/internal/routes"
	"github.com/go-chi/chi/v5"
	
)

func main() {
	// Load env variables
	config.LoadEnv()

	DB_URI := config.LoadMongoDBURI()
	DB_FOR_AUTH := config.LoadMongoDBNameAuth()
	DB_COLLECTION_FOR_AUTH := config.LoadMongoDBCollectionNameAuth()

	// Creating a MongoDB client using Db() function in db.go
	client := database.DbConnect(DB_URI)

	database.InitLoginCollection(client, DB_FOR_AUTH, DB_COLLECTION_FOR_AUTH)

	// Initialize Chi router
	r := chi.NewRouter()

	// Setting up middlewares
	
	// Invoking routes with the Chi router
	routes.Routes(r)

	// Defer disconnecting from the MongoDB client
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Panic("Error while disconnecting MongoDB client: ", err)
		}
	}()

	// Starting server
	SERVER_PORT := config.LoadServerPort()
	log.Printf("Server starting on port %s....", SERVER_PORT)
	err := http.ListenAndServe(":"+SERVER_PORT, r)
	if err != nil {
		log.Panic("Error while starting server: ", err)
	}
}
