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

func LoadMongoDBNameAuth() string{
	
	DB_NAME:=os.Getenv("DB_FOR_AUTH")
	if DB_NAME==""{
		log.Println("Error loading DB_FOR_AUTH in config.LoadMongoDBName()")
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