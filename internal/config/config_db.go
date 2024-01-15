package config

import (
	"log"
	"os"
)

func LoadMongoDBURI() string{
	
	DB_URI:=os.Getenv("MONGODB_URI")
	if DB_URI==""{
		log.Println("Error loading MONGODB_URI in config.LoadMongoDBURI()")
		return "mongodb://localhost:27017"
	}
	return DB_URI
}

func LoadMongoDBName() string{
	
	DB_NAME:=os.Getenv("DB_NAME")
	if DB_NAME==""{
		log.Println("Error loading DB_NAME in config.LoadMongoDBName()")
		return "jwt-auth-golang"
	}
	return DB_NAME
}


func LoadMongoDBCollectionNameAuth() string{
	
	DB_COLLECTION_NAME:=os.Getenv("DB_COLLECTION_FOR_AUTH")
	if DB_COLLECTION_NAME==""{
		log.Println("Error loading DB_COLLECTION_FOR_AUTH in config.LoadMongoDBCollectionName()")
		return "users"
	}
	return DB_COLLECTION_NAME
}

func LoadMongoDBCollectionEvent() string{
	
	DB_COLLECTION_NAME:=os.Getenv("DB_COLLECTION_FOR_EVENT")
	if DB_COLLECTION_NAME==""{
		log.Println("Error loading DB_COLLECTION_FOR_EVENT in config.LoadMongoDBCollectionName()")
		return "events"
	}
	return DB_COLLECTION_NAME
}