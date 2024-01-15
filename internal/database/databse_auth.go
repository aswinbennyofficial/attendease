package database

import (
	"errors"

	"context"
	"log"

	"github.com/aswinbennyofficial/attendease/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// has mongodb collection object for login
var coll *mongo.Collection

// InitLoginCollection initializes the mongodb collection for login
func InitLoginCollection(client *mongo.Client, dbName, collName string) error {
	// Initialisinbg nongodb collection object for login
	coll = client.Database(dbName).Collection(collName)
	return nil
}

// AddUserToDb adds a new user to the database
func AddUserToDb(newuser models.NewUser) error {
	result, err := coll.InsertOne(context.TODO(), newuser)
	if err != nil {
		return (err)
	}
	log.Println(result.InsertedID)
	return nil
}

// GetPasswordHashFromDb gets the password hash from the database
func GetHashAndUsernameFromDb(organisation string) (string,string, error) {
	// Checking if user exists
	isUserExist, err := DoesExist("organisation",organisation)
	if err != nil {
		log.Println("GetPasswordHashFromDb() :", err)
		return "","", err
	}
	if isUserExist == false {
		return "","", errors.New("Organisation does not exist")
	}

	// Creating a filter
	filter := bson.D{{"organisation", organisation}}

	// Instance of the NewUser struct
	var result models.NewUser

	// Find and decode from mongodb
	err = coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Println("GetPasswordHashFromDb() ", err)
		return "","", err
	}


	return result.Username,result.Password, nil
}

// DoesUserExist checks if a user exists in the database
func DoesExist(query string,username string) (bool, error) {
	opts := options.Count().SetHint("_id_")
	// Creating a filter
	filter := bson.D{{query, username}}
	// Counting the number of documents
	count, err := coll.CountDocuments(context.TODO(), filter, opts)
	if err != nil {
		log.Println(err)
		return true, err
	}
	if count == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

