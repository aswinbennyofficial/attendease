package controllers

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aswinbennyofficial/attendease/internal/config"
	"github.com/aswinbennyofficial/attendease/internal/database"
	"github.com/aswinbennyofficial/attendease/internal/models"
	"github.com/aswinbennyofficial/attendease/internal/utility"

	"github.com/xuri/excelize/v2"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
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
		doesEventIDexist,err= database.DoesExistInEventColl("eventid",event.EventId)
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

	for _,events:=range eventslist{
		eventListMap[events.EventId]=events
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	json.NewEncoder(w).Encode(eventListMap)
}

func HandleGetAnEvent(w http.ResponseWriter, r *http.Request) {
	// Get claims from context
	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		log.Println("Claims not found in context")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Claims not found in context"))
		return
	}

	eventid:=chi.URLParam(r, "eventid")
	if eventid==""{
		log.Println("EventID is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}


	event,err:=database.GetAnEventFromDb(claims.Org,eventid)

	if err!=nil{
		if err == mongo.ErrNoDocuments {
			log.Println("Event not found in database: ",err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Event not found in database"))
			return
		}
		log.Println("Error getting event from database: ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("Event fetched from database: ",event)
	
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	json.NewEncoder(w).Encode(event)
}

func HandleUploadParticipants(w http.ResponseWriter, r *http.Request){
	// Get claims from context
	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		log.Println("Claims not found in context")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Claims not found in context"))
		return
	}

	
	log.Println("Claims: ",claims.Org)

	// Get eventid from URL
	eventid:=chi.URLParam(r, "eventid")
	if eventid==""{
		log.Println("EventID is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if event exists in database
	doesEventIDexist,err:= database.DoesEventExistInSameOrg(claims.Org,eventid)
	if err!=nil{
		log.Println("Error checking if event exists in database: ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !doesEventIDexist{
		log.Println("Event does not exist in database")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Event does not exist in database"))
		return
	}
	
	// Parse multipart form
	err=r.ParseMultipartForm(10<<20)
	if err!=nil{
		log.Println("Error parsing multipart form: ",err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error parsing multipart form"))
		return
	}

	// Get file from form
	file,fileinfo,err:=r.FormFile("file")
	if err!=nil{
		log.Println("Error getting file from form: ",err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error getting file from form"))
		return
	}
	defer file.Close()
	log.Println("File uploaded: ",fileinfo.Filename)

	var participants []interface{}

	// Creating a map to prevent duplicate entries
	uidMap := make(map[string]bool)

	// Read and manipulate the CSV data
	csvReader := csv.NewReader(file)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			http.Error(w, "Error reading CSV file", http.StatusInternalServerError)
			return
		}
		
		
		var participant models.Participants

		// Generating a unique id
		uid,err:=utility.CreateRandomString(6)
		if err!=nil{
			log.Println("Error generating uid: ",err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Checking if the uid is already present in the map
		for uidMap[uid]{
			uid,err=utility.CreateRandomString(6)
			if err!=nil{
				log.Println("Error generating uid: ",err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		// adding the unique string to the map
		uidMap[uid]=true

		participant.ParticapantID=eventid+"-"+uid
		participant.Organisation=claims.Org
		participant.Name=record[0]
		participant.Email=record[1]
		participant.EventID=eventid
		participant.ScansCount=0

		participants=append(participants,participant)
	}

	log.Println("Participants list generated: ",participants)

	
	err=database.AddParticipantsToDb(participants)
	if err!=nil{
		log.Println("Error adding participants to database: ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("Participants list uploaded successfully")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Participants list uploaded successfully"))	
	
}

func HandleScan(w http.ResponseWriter, r *http.Request){
	// Get claims from context
	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		log.Println("Claims not found in context")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Claims not found in context"))
		return
	}

	var scanmodel models.ScanInput
	err:=json.NewDecoder(r.Body).Decode(&scanmodel)
	if err!=nil{
		log.Println("Error decoding scanmodel in HandleScan()")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	participantName,scancount,err:=database.AddScanToDb(claims.Org,scanmodel.ParticipantID,claims.Username)
	if err!=nil{
		if err== mongo.ErrNoDocuments{
			log.Println("Participant not found in database: ",err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Participant not found in database"))
			return
		}
		log.Println("Error adding scan to database: ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("Scan added to database of ", participantName, " ", scancount)

	var scanresponse models.ScanResponse
	scanresponse.Name=participantName
	scanresponse.ScansCount=scancount

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scanresponse)


}


func HandleGetParticipants(w http.ResponseWriter, r *http.Request){
	// Get claims from context
	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		log.Println("Claims not found in context")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Claims not found in context"))
		return
	}

	eventid:=chi.URLParam(r, "eventid")
	if eventid==""{
		log.Println("EventID is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	participants,err:=database.GetParticipantsFromDb(claims.Org,eventid)
	if err!=nil{
		log.Println("Error getting participants from database: ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("Participants fetched from database: ",participants)
	
	
	participantListMap := make(map[string]models.Participants)

	for _,participant:=range participants{
		participantListMap[participant.ParticapantID]=participant
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	json.NewEncoder(w).Encode(participantListMap)
}

func HandleGetParticipantsFile(w http.ResponseWriter, r *http.Request){
	// Get claims from context
	claims, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		log.Println("Claims not found in context")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Claims not found in context"))
		return
	}

	eventid:=chi.URLParam(r, "eventid")
	if eventid==""{
		log.Println("EventID is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	participants,err:=database.GetParticipantsFromDb(claims.Org,eventid)
	if err!=nil{
		log.Println("Error getting participants from database: ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("Participants fetched from database: ",participants)
	
	// Using excelize to create excel file
	f := excelize.NewFile()
    defer func() {
        if err := f.Close(); err != nil {
            log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return

        }
    }()
    // Create a new sheet.
    index, err := f.NewSheet("AllParticipants")
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
		return

    }
	// Create a new sheet.
	_, err = f.NewSheet("Present")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Create a new sheet.
	_, err = f.NewSheet("Absent")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set column headers.
	headers := map[string]string{
		"A1": "Name",
		"B1": "Email",
		"C1": "ParticipantID",
		"D1": "ScansCount",
	}

	for cell, value := range headers {
		f.SetCellValue("AllParticipants", cell, value)
		f.SetCellValue("Present", cell, value)
		f.SetCellValue("Absent", cell, value)
	}

	// Set active sheet of the workbook.
    f.SetActiveSheet(index)

	// Write all participant data to the sheet.
	for i, participant := range participants {
		rowNumber := strconv.Itoa(i + 2)
		f.SetCellValue("AllParticipants", "A"+rowNumber, participant.Name)
		f.SetCellValue("AllParticipants", "B"+rowNumber, participant.Email)
		f.SetCellValue("AllParticipants", "C"+rowNumber, participant.ParticapantID)
		f.SetCellValue("AllParticipants", "D"+rowNumber, participant.ScansCount)
	}
	// Write all  present participant data to the sheet.
	i:=0
	for _, participant := range participants {
		if participant.ScansCount>0{
		rowNumber := strconv.Itoa(i + 2)
		f.SetCellValue("Present", "A"+rowNumber, participant.Name)
		f.SetCellValue("Present", "B"+rowNumber, participant.Email)
		f.SetCellValue("Present", "C"+rowNumber, participant.ParticapantID)
		f.SetCellValue("Present", "D"+rowNumber, participant.ScansCount)
		i++
		}
	}

	// Write all  absent participant data to the sheet.
	i=0
	for _, participant := range participants {
		if participant.ScansCount==0{
		rowNumber := strconv.Itoa(i + 2)
		f.SetCellValue("Absent", "A"+rowNumber, participant.Name)
		f.SetCellValue("Absent", "B"+rowNumber, participant.Email)
		f.SetCellValue("Absent", "C"+rowNumber, participant.ParticapantID)
		f.SetCellValue("Absent", "D"+rowNumber, participant.ScansCount)
		i++
		}
	}

	defer func() {
		if err := os.Remove("AllParticipants.xlsx"); err != nil {
			log.Println("Error deleting XLSX file:", err)
		}
    }()

	// Save xlsx file by the given path.
	if err := f.SaveAs("AllParticipants.xlsx"); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}


	w.Header().Set("Content-Disposition", "attachment; filename=AllParticipants.xlsx")
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	http.ServeFile(w, r, "AllParticipants.xlsx")

	
}

func HandleSendEmail(w http.ResponseWriter, r *http.Request){
	// Get claims from context
	_, ok := r.Context().Value("claims").(*models.Claims)
	if !ok {
		log.Println("Claims not found in context")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Claims not found in context"))
		return
	}
	
	eventid:=chi.URLParam(r, "eventid")
	if eventid==""{
		log.Println("EventID is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err:=database.SendEmailToParticipants(eventid)
	if err!=nil{
		log.Println("Error sending email to participants: ",err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Email Queued successfully"))
}