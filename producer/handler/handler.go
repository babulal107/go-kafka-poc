package handler

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type MessagePayload struct {
	ID           string `json:"id"`
	MobileNumber string `json:"mobile_number"`
	Name         string `json:"name"`
	Message      string `json:"message"`
	Item         string `json:"item"`
}

func Order(c *gin.Context) {

	topic := "test-topic"

	request := MessagePayload{}
	err := c.ShouldBindBodyWithJSON(&request)
	if err != nil {
		fmt.Println("error bad request:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// Place order logic

	// Push a message to kafka queue
	// convert body into bytes and send it to kafka
	msgInBytes, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
	}

	if err = PushMessageToQueue(topic, msgInBytes); err != nil {
		// just log error no need to return error
		fmt.Println("error failed to push message in kafka queue:", err)
	}

	c.JSON(200, gin.H{
		"code":    http.StatusOK,
		"message": "Order placed",
	})
	return
}

func ConnectProducer(brokersUrl []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll // Ensure durability
	config.Producer.Retry.Max = 5                    // Retry sending messages
	config.Producer.Return.Successes = true          // Return successes

	producer, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func PushMessageToQueue(topic string, message []byte) error {
	brokers := []string{"localhost:29092"} // Change to your IBM Kafka broker address if needed
	producer, err := ConnectProducer(brokers)
	if err != nil {
		return err
	}
	defer producer.Close()

	// Message to send
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	// Send Message
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return err
	}

	fmt.Printf("Message sent successfully on topic=%s Partition=%d, Offset=%d\n", topic, partition, offset)

	return nil

}
