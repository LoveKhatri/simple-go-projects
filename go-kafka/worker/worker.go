package main

import (
	"fmt"       // For formatted I/O
	"os"        // For operating system functions
	"os/signal" // For signal handling
	"syscall"   // For system call constants

	"github.com/IBM/sarama" // Kafka client library
)

func main() {
	topic := "comments" // Kafka topic to consume messages from

	brokerUrl := []string{"localhost:9092"} // Kafka default runs on port 9092

	worker, err := connectConsumer(brokerUrl) // Connecting to Kafka consumer
	if err != nil {
		panic(err) // Panic if there is an error connecting to the consumer
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest) // Consuming messages from the oldest offset
	if err != nil {
		panic(err) // Panic if there is an error consuming the partition
	}

	fmt.Println("Consumer started") // Indicating that the consumer has started

	sigchan := make(chan os.Signal, 1) // Creating a channel to listen for OS signals

	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM) // Notifying the channel on interrupt signals

	msgCount := 0 // Counter for the number of messages processed

	doneChan := make(chan struct{}) // Channel to signal when processing is done

	go func() {
		for {
			select {
			case err := <-consumer.Errors(): // Handling consumer errors
				fmt.Println(err)

			case msg := <-consumer.Messages(): // Handling incoming messages
				msgCount++
				// Logging the received message details, similar to console.log in JavaScript
				fmt.Println("Received message", string(msg.Value), "| Topic: ", string(msg.Topic), "| Count:", msgCount)
			case <-sigchan: // Handling interrupt signals
				fmt.Println("Interrupt is detected")
				doneChan <- struct{}{}
			}
		}
	}()

	<-doneChan                                     // Waiting for the done signal
	fmt.Println("Processed", msgCount, "messages") // Logging the number of processed messages
	if err := worker.Close(); err != nil {         // Closing the consumer connection
		panic(err)
	}
}

func connectConsumer(brokerUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()         // Creating a new Sarama config
	config.Consumer.Return.Errors = true // Ensuring the consumer returns errors

	consumer, err := sarama.NewConsumer(brokerUrl, config) // Connecting to the Kafka consumer
	if err != nil {
		return nil, err
	}

	return consumer, nil // Returning the consumer connection
}
