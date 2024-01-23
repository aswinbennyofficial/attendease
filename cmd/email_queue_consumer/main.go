package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/skip2/go-qrcode"
)

var(
	// Create a channel to receive messages from the queue
	messageChannel chan amqp.Delivery=make(chan amqp.Delivery)

	RABBIT_MQ_URI string
	SMTP_USERNAME string
	SMTP_PASSWORD string
	SMTP_HOST string
	SMTP_PORT string
	FROM_EMAIL string
	REPLY_TO string

	// Number of workers
	workers int=10

)
func failOnError(err error, msg string) {
	if err != nil {
	  log.Panicf("%s: %s", msg, err)
	}
  }

func main() {
	// Load env variables
	err:=godotenv.Load()
	RABBIT_MQ_URI=os.Getenv("RABBIT_MQ_URI")
	SMTP_USERNAME=os.Getenv("SMTP_USERNAME")
	SMTP_PASSWORD=os.Getenv("SMTP_PASSWORD")
	SMTP_HOST =os.Getenv("SMTP_HOST")
	SMTP_PORT =os.Getenv("SMTP_PORT")
	FROM_EMAIL=os.Getenv("FROM_EMAIL")
	REPLY_TO =os.Getenv("REPLY_TO")


	log.Println("SMTP CREDS init ", SMTP_USERNAME, " ", SMTP_PASSWORD, " ", SMTP_HOST)

	// Setup authentication variable for email
	auth := smtp.PlainAuth("", SMTP_USERNAME, SMTP_PASSWORD, SMTP_HOST)
	for i:=0;i<workers;i++{
		go SendEmailRoutine(auth)
	}

	if REPLY_TO == "" {
		REPLY_TO = FROM_EMAIL
	}

	// Establish connection to RabbitMQ server
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

	// Consume messages from the queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	  )
	  failOnError(err, "Failed to register a consumer")
	  
	  // Create a channel to receive messages from the queue
	  var forever chan struct{}
	  
	 
	  
	  for d := range msgs {
		// Create a goroutine to consume messages from the queue
		messageChannel <- d
		}
		
	
	  
	  log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	 // Loop forever 
	  <-forever

}


 // Create a goroutine to consume messages from the queue
func SendEmailRoutine(auth smtp.Auth) {
	for messageDelivery := range messageChannel {
		
		log.Printf("Received a message: %s", messageDelivery.Body)

		// Split the message into substrings using the delimiter ":&:"
		substring := strings.Split(string(messageDelivery.Body), ":&:")
		if len(substring) != 8 {
			log.Println("Invalid queue element format", messageDelivery.Body)
			return
		}

		email := substring[0]
		name := substring[1]
		participantID := substring[2]
		eventName := substring[3]
		eventDate := substring[4]
		eventTime := substring[5]
		eventLocation := substring[6]
		eventDescription := substring[7]

		// Generate QR code
		qrCode, err := generateQRCode(participantID)
		if err != nil {
			log.Println("Error generating QR code:", err)
			return
		}

		// Convert QR code to PNG
		pngData, err := convertQRCodeToPNG(qrCode)
		if err != nil {
			log.Println("Error converting QR code to PNG:", err)
			return
		}

		// Create email body
		subject := "You are invited to " + eventName + "! 🎉"
		body := fmt.Sprintf("<html><body><h1>Hi %s,</h1> <p>You are invited to %s happening at %s %s on location %s</p> <p>%s</p> <p>this is an HTML-rich email template!<p><br><img src=\"cid:qrcode\"></body></html>", name, eventName, eventDate, eventTime, eventLocation, eventDescription)

		// Create MIME email with embedded image
		msg := []byte(
			"From: " + FROM_EMAIL + "\r\n" +
				"Reply-To: " + REPLY_TO + "\r\n" +
				"Subject: " + subject + "\r\n" +
				"MIME-version: 1.0;\nContent-Type: multipart/related; boundary=\"related_boundary\";\r\n" +
				"\r\n" +
				"--related_boundary\r\n" +
				"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
				"\r\n" +
				body + "\r\n" +
				"--related_boundary\r\n" +
				"Content-Type: image/png; name=\"qrcode.png\"\r\n" +
				"Content-Disposition: inline; filename=\"qrcode.png\"\r\n" +
				"Content-ID: <qrcode>\r\n" +
				"Content-Transfer-Encoding: base64\r\n" +
				"\r\n" +
				pngData + "\r\n" +
				"--related_boundary--")

		// Recipient email
		recieverEmail := []string{email}
		// Send email
		err = smtp.SendMail(SMTP_HOST+":"+SMTP_PORT, auth, FROM_EMAIL, recieverEmail, msg)

		// Handle errors
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Email sent successfully to ", email)
	}
}

// Generate QR code
func generateQRCode(data string) (*qrcode.QRCode, error) {
	qrCode, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	return qrCode, nil
}

// Convert QR code to PNG
func convertQRCodeToPNG(qrCode *qrcode.QRCode) (string, error) {
	// Create a new PNG image from the QR code
	pngData, err := qrCode.PNG(256)
	if err != nil {
		return "", err
	}

	// Encode the PNG image to base64
	pngBase64 := base64.StdEncoding.EncodeToString(pngData)

	return pngBase64, nil
}
