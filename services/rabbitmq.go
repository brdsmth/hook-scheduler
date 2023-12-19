package services

import (
	"context"
	"encoding/json"
	localConfig "hook-scheduler/config"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitMQConn *amqp.Connection
var RabbitMQConnMutex sync.Mutex

func ConnectToRabbitMQ(rabbitMQURL string) {
	var err error
	RabbitMQConn, err = amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	log.Println("Connected to RabbitMQ successfully")
}

func Queue(jobs []Job) {

	// Load RabbitMQ
	RabbitMQConnMutex.Lock()
	defer RabbitMQConnMutex.Unlock()

	if RabbitMQConn == nil {
		log.Println("RabbitMQ connection is not available")
		return
	}

	// Pull the name of the RabbitMQ queue from env
	rabbitMQQueue := localConfig.ReadEnv("RABBITMQ_QUEUE")
	if rabbitMQQueue == "" {
		log.Fatal("RABBITMQ_QUEUE environment variable not set")
	}

	// Create a channel
	ch, err := RabbitMQConn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return
	}
	defer ch.Close()

	// Declare the queue
	_, err = ch.QueueDeclare(rabbitMQQueue, false, false, false, false, nil)
	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return
	}

	// Use a background context or pass a context from higher up if available
	ctx := context.Background()

	// Publish to queue
	for _, job := range jobs {
		// Serialize job data
		jobData, err := json.Marshal(job)
		if err != nil {
			log.Printf("Error marshaling job: %v", err)
			continue // Skip this job and move to the next
		}

		err = ch.PublishWithContext(ctx, "", rabbitMQQueue, false, false, amqp.Publishing{
			// INFO -> below are additional details possible in RabbitMQ
			// Headers Table
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
			// Body []byte				 // Send the serialized job data
			ContentType: "application/json",
			Body:        jobData,
		})
		if err != nil {
			log.Printf("Failed to publish job to RabbitMQ: %v", err)
		} else {
			log.Printf("Published job:\t%s", job.RowKey)
		}
	}
}
