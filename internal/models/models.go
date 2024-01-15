package models

import(
	"time"
)


type Event struct {
	EventId string `json:"eventid" bson:"eventid"`
	EventName string `json:"eventname" bson:"eventname"`
	EventDescription string `json:"eventdescription" bson:"eventdescription"`
	EventLocation string `json:"eventlocation" bson:"eventlocation"`
	EventDate string `json:"eventdate" bson:"eventdate"`
	EventTime string `json:"eventtime" bson:"eventtime"`
	Organisation string `json:"organisation" bson:"organisation"`
	
}



type ScanInfo struct{
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	ScannedBy string `json:"scannedby"`
}

type Particpants struct{
	ParticapantID string `json:"particapantid"`
	Organisation string `json:"organisation"`
	Name string `json:"name"`
	Email string `json:"email"`
	EventID string `json:"eventid"`
	ScansCount int `json:"scancount"`
	ScansInfo []ScanInfo `json:"scansinfo"`

}

