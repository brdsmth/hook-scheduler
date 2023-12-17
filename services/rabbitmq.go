package services

import (
	"context"
	"log"
	"sync"

	gonanoid "github.com/matoous/go-nanoid"
	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitMQConn *amqp.Connection
var RabbitMQConnMutex sync.Mutex

type Job struct {
	Data string
}

func ConnectToRabbitMQ(rabbitMQURL string) {
	var err error
	RabbitMQConn, err = amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	log.Println("Connected to RabbitMQ successfully")
}

func Queue(jobs []*Job) {

	// Load RabbitMQ
	RabbitMQConnMutex.Lock()
	defer RabbitMQConnMutex.Unlock()

	if RabbitMQConn == nil {
		log.Println("RabbitMQ connection is not available")
		return
	}

	// Create a channel
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return
	}
	defer ch.Close()

	// Declare the queue
	_, err = ch.QueueDeclare("QUEUE", false, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return
	}

	// Use a background context or pass a context from higher up if available
	ctx := context.Background()

	// Publish to queue
	for _, job := range jobs {
		// Identify job
		id, err := gonanoid.Nanoid(6)
		if err != nil {
			log.Printf("error adding nanoid: %s", err)
			return
		}
		log.Printf("Queue job:\t\t%s", id)

		err = ch.PublishWithContext(ctx, "", "QUEUE", false, false, amqp.Publishing{
			// // Application or exchange specific fields,
			// // the headers exchange will inspect this field.
			// Headers Table
			// // Properties
			// ContentType     string    // MIME content type
			// ContentEncoding string    // MIME content encoding
			// DeliveryMode    uint8     // Transient (0 or 1) or Persistent (2)
			// Priority        uint8     // 0 to 9
			// CorrelationId   string    // correlation identifier
			// ReplyTo         string    // address to to reply to (ex: RPC)
			// Expiration      string    // message expiration spec
			// MessageId       string    // message identifier
			// Timestamp       time.Time // message timestamp
			// Type            string    // message type name
			// UserId          string    // creating user id - ex: "guest"
			// AppId           string    // creating application id

			// // The application specific payload of the message
			// Body []byte
			ContentType: "text/plain",
			Body:        []byte(job.Data), // Assuming Task has a Data field
		})
		if err != nil {
			log.Printf("Failed to publish task: %v", err)
		}
	}
}
