package main

import (
	"github.com/babulal107/go-kafka-poc/consumer/service"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const workerCount = 5 // Number of concurrent workers

func main() {
	brokers := []string{"localhost:29092"}
	topic := "test-topic"

	// Create a buffered channel for processing messages
	messageChannel := make(chan service.KafkaMessagePayload, 100) // Buffer size 100

	// Asynchronous with Worker Pool → Run multiple goroutines consuming from the channel for higher throughput.
	// Start the worker pool to read a message from the channel and call to send SMS service
	service.StartWorkerPool(workerCount, messageChannel)

	// Asynchronous with Channel only → Using a Channel for Asynchronous SMS Processing
	// Start SMS processing
	//go service.ProcessMessagesFromChannel(messageChannel)

	// Start Kafka consumer
	consumerService, err := service.NewKafkaConsumer(brokers, topic, messageChannel)
	if err != nil {
		log.Fatalf("Failed to initialize Kafka consumer: %v\n", err)
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	go consumerService.ConsumeMessages()

	<-stopChan
	log.Println("Received termination signal. Shutting down gracefully...")

	consumerService.Close()
	close(messageChannel) // Close the channel to stop workers

	log.Println("Consumer shut down successfully.")
}
