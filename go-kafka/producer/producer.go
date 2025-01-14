package main

import (
	"encoding/json" // For encoding and decoding JSON
	"fmt"           // For formatted I/O
	"log"           // For logging

	"github.com/IBM/sarama"       // Kafka client library
	"github.com/gofiber/fiber/v2" // Web framework similar to Express
)

// Comment struct to map the incoming JSON data
type Comment struct {
	Text string `form:"text" json:"text"` // Similar to defining a schema in Mongoose
}

func main() {
	app := fiber.New()          // Initializing a new Fiber app, similar to Express
	api := app.Group("/api/v1") // Grouping routes under /api/v1

	api.Post("/comment", createComment) // Defining a POST route for /comment

	app.Listen(":3000") // Starting the server on port 3000
}

func createComment(c *fiber.Ctx) error {
	cmt := new(Comment) // Creating a new instance of Comment struct

	if err := c.BodyParser(cmt); err != nil { // Parsing the request body into cmt
		log.Println(err) // Logging the error

		// Returning a 400 error response, similar to res.status(400).json({ success: false, message: err }) in Express
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})

		return err
	}

	cmtInBytes, err := json.Marshal(cmt) // Converting the Comment struct to JSON
	if err != nil {
		log.Println(err)
		// Returning a 500 error response, similar to res.status(500).json({ success: false, message: err }) in Express
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
	}

	PushCommentToQueue("comments", cmtInBytes) // Pushing the comment to Kafka queue

	// Returning a success response, similar to res.json({ success: true, message: "Comment pushed successfully", comment: cmt }) in Express
	err = c.JSON(&fiber.Map{
		"success": true,
		"message": "Comment pushed successfully",
		"comment": cmt,
	})
	if err != nil {
		log.Println(err)
		// Returning a 500 error response
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})

		return err
	}

	return err
}

func PushCommentToQueue(topic string, message []byte) error {
	brokerUrl := []string{"localhost:9092"}     // Kafka default runs on port 9092
	producer, err := ConnectProducer(brokerUrl) // Connecting to Kafka producer
	if err != nil {
		log.Println(err)
		return err
	}

	defer producer.Close() // Ensuring the producer connection is closed

	msg := sarama.ProducerMessage{ // Creating a new Kafka message
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(&msg) // Sending the message to Kafka
	if err != nil {
		log.Println(err)
		return err
	}

	// Logging the message details, similar to console.log in JavaScript
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)

	return nil
}

func ConnectProducer(brokerUrl []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()                     // Creating a new Sarama config
	config.Producer.Return.Successes = true          // Ensuring the producer returns successes
	config.Producer.RequiredAcks = sarama.WaitForAll // Waiting for all in-sync replicas to acknowledge
	config.Producer.Retry.Max = 5                    // Setting the retry count

	conn, err := sarama.NewSyncProducer(brokerUrl, config) // Connecting to the Kafka producer
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}
