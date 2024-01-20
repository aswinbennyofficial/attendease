package utility

import (
	"context"
	"log"
	"time"
	

	"github.com/aswinbennyofficial/attendease/internal/config"
	"github.com/aswinbennyofficial/attendease/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
	
)

// Helper function to check the return value for each amqp call:
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}


func EmailQueueSender(participants []models.Participants,event models.Event) error{
	RABBIT_MQ_URI:=config.LoadRabbitMQURI()

	

	// Connect to RabbitMQ server
	conn, err := amqp.Dial(RABBIT_MQ_URI)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declaring a queue
	q, err := ch.QueueDeclare(
		"emailchan", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, participant := range participants{ 
		// Publish a message to the queue
		body := participant.Email+":&:"+participant.Name+":&:"+participant.ParticapantID+":&:"+event.EventName+":&:"+event.EventDate+":&:"+event.EventTime+":&:"+event.EventLocation+":&:"+event.EventDescription

		err = ch.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] queued %s\n", body)

		
	}
	return nil
}