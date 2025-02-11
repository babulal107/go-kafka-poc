package service

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

// KafkaConsumer struct handles Kafka consumption
type KafkaConsumer struct {
	consumer       sarama.PartitionConsumer
	topic          string
	messageChannel chan KafkaMessagePayload
}

// NewKafkaConsumer initializes a new Kafka consumer
func NewKafkaConsumer(brokers []string, topic string, messageChannel chan KafkaMessagePayload) (*KafkaConsumer, error) {
	consumer, err := connectConsumer(brokers, topic)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer:       consumer,
		topic:          topic,
		messageChannel: messageChannel,
	}, nil
}

// ConsumeMessages starts consuming messages
func (kc *KafkaConsumer) ConsumeMessages() {
	log.Printf("Waiting for messages on topic: %s", kc.topic)
	msgCount := 0
	for {
		select {
		case msg := <-kc.consumer.Messages():
			msgCount++
			fmt.Printf("Received message Count: %d | Topic: (%s) | Message: %s\n", msgCount, msg.Topic, string(msg.Value))

			payload, err := validateAndExtractMessage(msg.Value)
			if err != nil {
				log.Printf("Invalid message: %v\n", err)
				continue
			}
			// Send the validated message to the channel for processing
			kc.messageChannel <- payload

		case err := <-kc.consumer.Errors():
			log.Printf("Consumer error: %v", err)
		}
	}
}

// Close shuts down the consumer
func (kc *KafkaConsumer) Close() {
	kc.consumer.Close()
}

// connectConsumer connects to Kafka
func connectConsumer(brokers []string, topic string) (sarama.PartitionConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	partConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return nil, err
	}

	return partConsumer, nil
}
