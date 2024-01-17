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
var employeecoll *mongo.Collection

// InitLoginCollection initializes the mongodb collection for login
func InitLoginCollection(client *mongo.Client, dbName, collName string) error {
	// Initialisinbg nongodb collection object for login
	coll = client.Database(dbName).Collection(collName)
	return nil
}

func InitEmployeeCollection(client *mongo.Client, dbName, collName string) error {
	// Initialisinbg nongodb collection object for login
	employeecoll = client.Database(dbName).Collection(collName)
	return nil
}

// AddUserToDb adds a new admin to the database
func AddAdminToDb(newuser models.NewUser) error {
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
	isUserExist, err := DoesExistInAuthColl("organisation",organisation)
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
func DoesExistInAuthColl(query string,value string) (bool, error) {
	opts := options.Count().SetHint("_id_")
	// Creating a filter
	filter := bson.D{{query, value}}
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

// Check if Employee exist in the same table
func DoesEmpExist(org string,username string) (bool, error) {
	opts := options.Count().SetHint("_id_")
	// Creating a filter
	filter := bson.D{{"$and", bson.A{     bson.D{{"organisation", org}},     bson.D{{"username", username}}, }}}
	// Counting the number of documents
	count, err := employeecoll.CountDocuments(context.TODO(), filter, opts)
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


func AddEmployeeToDb(newuser models.NewUser) error {
	result, err := employeecoll.InsertOne(context.TODO(), newuser)
	if err != nil {
		return (err)
	}
	log.Println(result.InsertedID)
	return nil
}

func GetEmployeesFromDb(organisation string) ([]string, error) {
	// Creating a filter
	filter := bson.D{{"organisation", organisation}}

	// Instance of the NewUser struct
	var results []string

	// Find and decode from mongodb
	cur, err := employeecoll.Find(context.Background(), filter)
	if err != nil {
		log.Println("GetEmployeesFromDb() ", err)
		return results, err
	}

	// Iterating through the results
	for cur.Next(context.Background()) {
		var result models.NewUser
		err := cur.Decode(&result)
		if err != nil {
			log.Println("GetEmployeesFromDb() ", err)
			return results, err
		}
		results = append(results, result.Username)
	}

	return results, nil
}

func GetHashFromEmployeeColl(organisation string,username string) (string, error) {
	// Checking if user exists
	isUserExist, err := DoesEmpExist(organisation,username)
	if err != nil {
		log.Println("GetHashFromEmployeeColl() :", err)
		return "", err
	}
	if isUserExist == false {
		return "", errors.New("Employee does not exist")
	}

	// Creating a filter
	filter := bson.D{{"$and", bson.A{     bson.D{{"organisation", organisation}},     bson.D{{"username", username}}, }}}

	// Instance of the NewUser struct
	var result models.NewUser

	// Find and decode from mongodb
	err = employeecoll.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		log.Println("GetPasswordHashFromDb() ", err)
		return "", err
	}
	return result.Password, nil
}