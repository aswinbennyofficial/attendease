package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aswinbennyofficial/attendease/internal/config"
	"github.com/aswinbennyofficial/attendease/internal/database"
	"github.com/aswinbennyofficial/attendease/internal/models"
	"github.com/aswinbennyofficial/attendease/internal/utility"
)

func HandleCreateEvent(w http.ResponseWriter, r *http.Request) {
	// Get claims from context
	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		log.Println("Claims not found in context")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Claims not found in context"))
		return
	}

	EventCollectionName:=config.LoadMongoDBCollectionEvent()
	log.Println("EventCollectionName fetched: ",EventCollectionName)

	var event models.Event
	err:=json.NewDecoder(r.Body).Decode(&event)
	if err!=nil{
		log.Println("Error decoding event in HandleCreateEvent()")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if event.EventName==""{
		log.Println("Event name is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	doesEventIDexist:=true
	for doesEventIDexist{
		event.EventId,err=utility.CreateRandomString(6)
		doesEventIDexist,err= database.DoesExistInAuthColl("eventid",event.EventId)
	}
	log.Println("EventID generated: ",event.EventId)

	event.Organisation=claims.Org
	err=database.AddEventToDb(event)
	if err!=nil{
		log.Println("Error adding event to database: ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return 			
	}
	log.Println("Event added to database")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Event created successfully"))
}

func HandleGetEvents(w http.ResponseWriter, r *http.Request) {
	// Get claims from context
	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		log.Println("Claims not found in context")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Claims not found in context"))
		return
	}

	eventslist,err:=database.GetEventsFromDb(claims.Org)
	if err!=nil{
		log.Println("Error getting events from database: ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("Events fetched from database: ",eventslist)
	
	
	eventListMap := make(map[string]models.Event)

	for _,event:=range eventslist{
		eventListMap[event.EventId]=event
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	json.NewEncoder(w).Encode(eventListMap)
}