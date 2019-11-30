package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// User profile structure
type User struct {
	Name   string `json:"Name"`
	Type   string `json:"Type"`
	Status string `json:"Status"`
}

var (
	// ErrNameNotProvided is thrown when a name is not provided
	ErrNameNotProvided = errors.New("no name was provided in the HTTP body")
	// ErrInvalidBody is thrown when a the body is not parsable
	ErrInvalidBody = errors.New("Request Body is invalid")
)

// Handler for  Amazon API Gateway request/responses
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var profile *User

	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)

	// If no name is provided in the HTTP request body, throw an error
	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
	}

	unmErr := json.Unmarshal([]byte(request.Body), &profile)
	if unmErr != nil {
		return events.APIGatewayProxyResponse{}, ErrInvalidBody
	}
	sendSqs(profile)

	return events.APIGatewayProxyResponse{
		Body:       "Hello " + string(profile.Name),
		StatusCode: 200,
	}, nil

}

func sendSqs(profile *User) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)
	qURL := os.Getenv("TASK_QUEUE_URL")

	result, err := svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Name": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(profile.Name),
			},
			"Type": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(profile.Type),
			},
			"Status": &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(profile.Status),
			},
		},
		MessageBody: aws.String("Information about user profile status"),
		QueueUrl:    &qURL,
	})

	if err != nil {
		log.Println("Error sending SQS message", err)
	} else {
		log.Println("Success", *result.MessageId)
	}

}

func main() {
	lambda.Start(Handler)
}
