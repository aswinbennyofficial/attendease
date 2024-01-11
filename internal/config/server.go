
package config

import (
	"os"
	"log"
)

func LoadServerPort() string{
	PORT:=os.Getenv("PORT")
	if PORT==""{
		log.Println("Error loading PORT in config.LoadPort()")
		return "8080"
	}
	return PORT
}