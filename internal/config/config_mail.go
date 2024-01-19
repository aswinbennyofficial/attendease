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

// func LoadSMTPUsername() string{
// 	SMTP_USERNAME:=os.Getenv("SMTP_USERNAME")
// 	if SMTP_USERNAME==""{
// 		log.Println("Errror loading SMTP_USERNAME in config.LoadRabbitMQURI()")
// 	}
// 	return SMTP_USERNAME
// }

// func LoadSMTPPassword() string{
// 	SMTP_PASSWORD:=os.Getenv("SMTP_PASSWORD")
// 	if SMTP_PASSWORD==""{
// 		log.Println("Errror loading SMTP_PASSWORD in config.LoadRabbitMQURI()")
// 	}
// 	return SMTP_PASSWORD
// }

// func LoadSMTPHost() string{
// 	SMTP_HOST:=os.Getenv("SMTP_HOST")
// 	if SMTP_HOST==""{
// 		log.Println("Errror loading SMTP_HOST in config.LoadRabbitMQURI()")
// 	}
// 	return SMTP_HOST
// }

// func LoadFromAddress() string{
// 	FROM_EMAIL:=os.Getenv("FROM_EMAIL")
// 	if FROM_EMAIL==""{
// 		log.Println("Errror loading FROM_EMAIL in config.LoadRabbitMQURI()")
// 	}
// 	return FROM_EMAIL
// }

// func LoadSMTPPort() string{
// 	SMTP_PORT:=os.Getenv("SMTP_PORT")
// 	if SMTP_PORT==""{
// 		log.Println("Errror loading SMTP_PORT in config.LoadRabbitMQURI()")
// 	}
// 	return SMTP_PORT
// }


// func LoadReplyToAddress() string{
// 	REPLY_TO:=os.Getenv("REPLY_TO")
// 	if REPLY_TO==""{
// 		log.Println("Errror loading REPLY_TO in config.LoadRabbitMQURI()")
// 	}
// 	return REPLY_TO
// }