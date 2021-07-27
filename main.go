package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// MarshalMap requires public struct fields
// if your key/attributes shhould be lowercase, add directives
type Record struct {
	Timestamp string `json:"timestamp"`
	Data      string `json:"data"`
}

func handleRequest(ctx context.Context, event events.SQSEvent) (string, error) {
	// event
	eventJson, _ := json.MarshalIndent(event, "", "  ")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Fatalf("Failed to start a session, %s", err)
		return "error", err
	}

	// // Create an SQS client
	// sqsClient := sqs.New(sess)

	// Create a DynamoDB client
	dynamodbClient := dynamodb.New(sess)
	for _, record := range event.Records {
		r := Record{
			Timestamp: time.Now().Format("20060102150405"),
			Data:      record.Body,
		}
		log.Printf("Writing %s to DynamoDB table", record.Body)
		itemInput, err := dynamodbattribute.MarshalMap(r)
		if err != nil {
			log.Fatalf("Failed to convert input data to a DynamoDB item format, %s", err)
		}

		_, err = dynamodbClient.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String("test"),
			Item:      itemInput,
		})
		if err != nil {
			log.Fatalf("Failed to PutItem: %s", err)
		}
	}
	return string(eventJson), nil
}

func main() {
	runtime.Start(handleRequest)
}
