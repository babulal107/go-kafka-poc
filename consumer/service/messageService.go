package service

import (
	"encoding/json"
	"fmt"
	"log"
)

type KafkaMessagePayload struct {
	ID           string `json:"id"`
	MobileNumber string `json:"mobile_number"`
	Name         string `json:"name"`
	Message      string `json:"message"`
	Item         string `json:"item"`
}

// ProcessMessagesFromChannel continuously processes messages from the channel
func ProcessMessagesFromChannel(messageChannel chan KafkaMessagePayload) {
	for payload := range messageChannel {
		// prepare send SMSPayload
		smsPayload := SMSPayload{
			SourceNumber:      "AMNOTIF",
			DestinationNumber: payload.MobileNumber,
			Message:           payload.Message,
		}

		if err := sendSMSWithRetries(smsPayload, 3); err != nil {
			log.Printf("failed to send SMS after retries: %v", err)
		} else {
			log.Printf("SMS sent successfully to %s", payload.MobileNumber)
		}
	}
}

// StartWorkerPool initializes multiple workers to process messages concurrently
func StartWorkerPool(workerCount int, messageChannel chan KafkaMessagePayload) {
	for i := 1; i <= workerCount; i++ {
		go worker(i, messageChannel)
	}
	log.Printf("%d workers started.", workerCount)
}

// worker function that processes messages from the channel
func worker(workerID int, messageChannel chan KafkaMessagePayload) {
	log.Printf("Worker %d started", workerID)

	for payload := range messageChannel {
		log.Printf("Worker %d processing message for %s", workerID, payload.MobileNumber)

		// prepare send SMSPayload
		smsPayload := SMSPayload{
			SourceNumber:      "AMNOTIF",
			DestinationNumber: payload.MobileNumber,
			Message:           payload.Message,
		}

		if err := sendSMSWithRetries(smsPayload, 3); err != nil {
			log.Printf("Worker %d failed to send SMS after retries: %v", workerID, err)
		} else {
			log.Printf("Worker %d successfully sent SMS to %s", workerID, payload.MobileNumber)
		}
	}
	log.Printf("Worker %d stopped", workerID)
}

// validateAndExtractMessage validates the Kafka message and extracts necessary fields
func validateAndExtractMessage(msgValue []byte) (KafkaMessagePayload, error) {
	var msgData KafkaMessagePayload
	if err := json.Unmarshal(msgValue, &msgData); err != nil {
		return msgData, fmt.Errorf("message is not valid JSON")
	}

	// TODO: do more validation if require

	return msgData, nil
}
