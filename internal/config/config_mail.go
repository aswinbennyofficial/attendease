package config

import (
	"os"
	"log"
)

func LoadRabbitMQURI() string{
	
	RABBITMQ_URI:=os.Getenv("RABBIT_MQ_URI")
	if RABBITMQ_URI==""{
		log.Println("Error loading RABBITMQ_URI in config.LoadRabbitMQURI()")
		return "amqp://guest:guest@localhost:5672/"
	}
	return RABBITMQ_URI
}