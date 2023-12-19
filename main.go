package main

import (
	"hook-scheduler/config"
	"hook-scheduler/services"
	"log"
	"net/http"
)

func main() {
	var err error

	// Connect to DynamoDB
	services.ConnectDynamoDB()

	// Connect to RabbitMQ
	// rabbitMQURL := "amqp://guest:guest@localhost:5672/" // local rabbitmq url
	rabbitMQURL := config.ReadEnv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		log.Fatal("RABBITMQ_URL environment variable not set")
	}

	services.ConnectToRabbitMQ(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer services.RabbitMQConn.Close()

	// Initiate cron
	services.RunCron()

	// Start the HTTP server for the publisher microservice
	log.Println("Server listening on 8082")
	http.ListenAndServe(":8082", nil)
}
