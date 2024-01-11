package database


import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"log"
)

func DbConnect(DB_URI string) *mongo.Client {

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	
	// Configue settings related Mongodb client behaviour
	opts := options.Client().ApplyURI(DB_URI).SetServerAPIOptions(serverAPI)
	
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	
	
	// Check the connection
	err = client.Ping(context.TODO(), nil) 
	if err != nil {
	log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")
	return client
}