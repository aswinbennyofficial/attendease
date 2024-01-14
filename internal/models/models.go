package models

import(
	"time"
)

type Events struct{
	EventID string `json:"eventid"`
	Organisation string `json:"organisation"`
	Name string `json:"name"`
	Description string `json:"description"`
	EventTime string `json:"eventtime"`
	Location string `json:"location"`
}

type ScanInfo struct{
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
	ScannedBy string `json:"scannedby"`
}

type Particapants struct{
	ParticapantID string `json:"particapantid"`
	Organisation string `json:"organisation"`
	Name string `json:"name"`
	Email string `json:"email"`
	EventID string `json:"eventid"`
	ScansCount int `json:"scancount"`
	ScansInfo []ScanInfo `json:"scansinfo"`

}

