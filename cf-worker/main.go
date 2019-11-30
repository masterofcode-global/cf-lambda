package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is your Lambda function handler
func Handler(ctx context.Context, sqsEvent events.SQSEvent) (string, error) {
	var name, profileType, status string
	name = ""
	profileType = ""
	status = ""

	for _, message := range sqsEvent.Records {
		log.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)

		if message.MessageAttributes["Name"].StringValue != nil {
			name = *(message.MessageAttributes["Name"].StringValue)
		}
		if message.MessageAttributes["Type"].StringValue != nil {
			profileType = *(message.MessageAttributes["Type"].StringValue)
		}
		if message.MessageAttributes["Status"].StringValue != nil {
			status = *(message.MessageAttributes["Status"].StringValue)
		}

		log.Printf("Name: %v, Type: %v, Status: %v \n", name, profileType, status)
	}

	return "Procesed messages", nil
}

func main() {
	lambda.Start(Handler)
}
