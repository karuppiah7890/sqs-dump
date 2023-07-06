package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/karuppiah7890/sqs-dump/pkg/config"

	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// TODO: Write tests for all of this

var version string

func checkSignal(signals chan os.Signal, done chan bool) {
	<-signals
	done <- true
}

func main() {
	done := make(chan bool, 1)
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, os.Interrupt)

	go checkSignal(signals, done)

	log.Printf("version: %v", version)
	c, err := config.NewConfigFromEnvVars()
	if err != nil {
		log.Fatalf("error occurred while getting configuration from environment variables: %v", err)
	}

	awsconfig, err := awsconf.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("error occurred while loading aws configuration: %v", err)
	}

	sqsClient := sqs.NewFromConfig(awsconfig)

	queueUrl := c.GetSqsQueueUrl()

	f, err := os.OpenFile("messages.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer closeFile(f)

	for {
		select {
		case <-done:
			os.Exit(0)

		default:
			messages, err := getMessagesFromQueue(queueUrl, sqsClient)
			if err != nil {
				log.Fatalf("error occurred while getting messages from sqs queue: %v", err)
			}

			for _, message := range messages {
				messageJson, err := json.Marshal(message)
				if err != nil {
					log.Fatalf("error occurred while getting serializing message from sqs queue: %v", err)
				}

				if _, err = f.WriteString(string(messageJson) + "\n"); err != nil {
					log.Fatalf("error occurred while writing message from sqs queue into file: %v", err)
				}

				fmt.Printf(".")
			}
		}

	}
}

type Message struct {
	// The message's contents (not URL-encoded).
	Body string `json:"body"`

	// An MD5 digest of the non-URL-encoded message body string.
	MD5OfBody string `json:"md5_of_body"`

	// A unique identifier for the message. A MessageIdis considered unique across all
	// Amazon Web Services accounts for an extended period of time.
	MessageId string `json:"message_id"`

	// An identifier associated with the act of receiving the message. A new receipt
	// handle is returned every time you receive a message. When deleting a message,
	// you provide the last received receipt handle to delete the message.
	ReceiptHandle string `json:"receipt_handle"`

	// Each message attribute consists of a Name, Type, and Value. For more
	// information, see Amazon SQS message attributes
	// (https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-message-metadata.html#sqs-message-attributes)
	// in the Amazon SQS Developer Guide.
	MessageAttributes map[string]types.MessageAttributeValue `json:"message_attributes"`
}

// Get message from the queue
func getMessagesFromQueue(queueUrl string, sqsClient *sqs.Client) ([]*Message, error) {
	input := sqs.ReceiveMessageInput{
		QueueUrl:              &queueUrl,
		MaxNumberOfMessages:   10,
		MessageAttributeNames: []string{"All"},
	}

	output, err := sqsClient.ReceiveMessage(context.TODO(), &input)
	if err != nil {
		return nil, fmt.Errorf("error occurred while receiving sqs queue message: %v", err)
	}

	messages := make([]*Message, 0)

	for _, outputMessage := range output.Messages {
		message := Message{
			Body:              *outputMessage.Body,
			MD5OfBody:         *outputMessage.MD5OfBody,
			MessageId:         *outputMessage.MessageId,
			ReceiptHandle:     *outputMessage.ReceiptHandle,
			MessageAttributes: outputMessage.MessageAttributes,
		}

		messages = append(messages, &message)
	}

	return messages, nil
}

func closeFile(file *os.File) {
	err := file.Close()
	log.Fatalf("error occurred while closing file: %v", err)
}
