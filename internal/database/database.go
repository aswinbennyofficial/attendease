package database

import (
	"context"
	"log"
	"time"

	"github.com/aswinbennyofficial/attendease/internal/models"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"
)

// has mongodb collection object for login
var Eventcoll *mongo.Collection

// has mongodb collection object for partcipant list
var Participantcoll *mongo.Collection

// InitLoginCollection initializes the mongodb collection for login
func InitEventCollection(client *mongo.Client, dbName, collName string) error {
	// Initialisinbg nongodb collection object for login
	Eventcoll = client.Database(dbName).Collection(collName)
	return nil
}

func InitParticipantCollection(client *mongo.Client, dbName, collName string) error {
	// Initialisinbg nongodb collection object for login
	Participantcoll = client.Database(dbName).Collection(collName)
	return nil
}

func AddEventToDb(event models.Event) error {
	// Inserting the event into the database
	_, err := Eventcoll.InsertOne(context.Background(), event)
	if err != nil {
		log.Println("Error while inserting event into database: ", err)
		return err
	}
	return nil
}

func GetEventsFromDb(organisation string) ([]models.Event, error) {
	// Getting all events from the database
	filter := bson.D{{"organisation", organisation}}
	cursor, err := Eventcoll.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error while getting events from database: ", err)
		return nil, err
	}
	defer cursor.Close(context.Background())
	var events []models.Event
	for cursor.Next(context.Background()) {
		var event models.Event
		cursor.Decode(&event)
		events = append(events, event)
	}
	return events, nil
}

func GetAnEventFromDb(organisation string,eventid string) (models.Event, error) {
	// Getting an event from the database
	//filter := bson.D{{"eventid", eventid} , {"organisation", organisation}}
	filter :=bson.D{{"$and", bson.A{     bson.D{{"eventid", eventid}},     bson.D{{"organisation", organisation}}, }}}

	var event models.Event
	err := Eventcoll.FindOne(context.Background(), filter).Decode(&event)
	if err != nil {
		log.Println("Error while getting event from database: ", err)
		return event, err
	}
	return event, nil
}


func AddParticipantsToDb(participants []interface{}) error {
	// Inserting the participant into the database
	_, err := Participantcoll.InsertMany(context.Background(), participants)
	if err != nil {
		log.Println("Error while inserting participant into database: ", err)
		return err
	}
	return nil
}

func AddScanToDb(organisation string,participantId string, employee string) (string,int,error){

	
    filter := bson.D{{"$and", bson.A{     bson.D{{"organisation", organisation}},     bson.D{{"particapantid", participantId}}, }}}

    // Define the update operation
    update := bson.D{
        {"$inc", bson.D{{"scanscount", 1}}},
        {"$push", bson.D{{"scansinfo", models.ScanInfo{Timestamp: time.Now(), ScannedBy: employee}}}},
    }

    // Perform the update
    _, err := Participantcoll.UpdateOne(context.Background(), filter, update)
    if err != nil {
        log.Println("Error while updating participant in database: ", err)
        return "", 0, err
    }

	 // Retrieve the updated participant
	 var updatedParticipant models.Participants
	 err = Participantcoll.FindOne(context.Background(), filter).Decode(&updatedParticipant)
	 if err != nil {
		 log.Println("Error while getting updated participant from database: ", err)
		 return "", 0, err
	 }
	return updatedParticipant.Name,updatedParticipant.ScansCount,nil

}