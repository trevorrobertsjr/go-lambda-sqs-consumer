package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	runtime "github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func handleRequest(ctx context.Context, event events.CodePipelineEvent) (string, error) {
	// event
	eventJson, _ := json.MarshalIndent(event, "", "  ")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Println("Failed to create bucket", err)
		return "error", err
	}

	// Create an SQS client
	sqsClient := sqs.New(sess)

	// Create a DynamoDB client
	dynamodbClient := dynamodb.New(sess)

	return string(eventJson), nil
}

func main() {
	runtime.Start(handleRequest)
}
