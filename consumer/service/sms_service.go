package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// SMSServiceURL - Replace this with the actual SMS API URL
const SMSServiceURL = "https://your-sms-service.com/api/send"

// SMSPayload defines the structure for the SMS API request
type SMSPayload struct {
	SourceNumber      string `json:"sourceNumber"`
	DestinationNumber string `json:"destinationNumber"`
	Message           string `json:"message"`
}

// sendSMSWithRetries attempts to send an SMS with retries and exponential backoff
func sendSMSWithRetries(payload SMSPayload, maxRetries int) error {
	var err error
	backoff := time.Second // Initial backoff delay

	for i := 0; i < maxRetries; i++ {
		err = sendSMS(payload)
		if err == nil {
			return nil // Success, exit retry loop
		}

		log.Printf("Retry %d/%d for %s failed: %v", i+1, maxRetries, payload.DestinationNumber, err)
		time.Sleep(backoff) // Exponential backoff
		backoff *= 2
	}

	return fmt.Errorf("failed to send SMS to %s after %d retries: %w", payload.DestinationNumber, maxRetries, err)
}

// sendSMS makes an HTTP POST request to the SMS service
func sendSMS(payload SMSPayload) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	fmt.Println("Send SMS Request : ", string(jsonData))

	req, err := http.NewRequest("POST", SMSServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SMS service responded with status: %d", resp.StatusCode)
	}

	return nil
}
