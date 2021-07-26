package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Record struct {
	timestamp string
	data      string
}

func handleRequest(ctx context.Context, event events.SQSEvent) (string, error) {
	// event
	eventJson, _ := json.MarshalIndent(event, "", "  ")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Println("Failed to get event data", err)
		return "error", err
	}

	// // Create an SQS client
	// sqsClient := sqs.New(sess)

	// Create a DynamoDB client
	dynamodbClient := dynamodb.New(sess)
	for _, record := range event.Records {
		r := Record{
			timestamp: time.Now().Format("20060102150405"),
			data:      record.Body,
		}
		av, err := dynamodbattribute.MarshalMap(r)
		if err != nil {
			panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
		}

		dynamodbClient.PutItem(&dynamodb.PutItemInput{

			TableName: aws.String("test"),
			Item:      av,
		})
	}
	return string(eventJson), nil
}

func main() {
	runtime.Start(handleRequest)
}
