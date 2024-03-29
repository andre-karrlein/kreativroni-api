package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/andre-karrlein/kreativroni-api/util"
)

// ErrNoCustomers indicates that API failed in some way
var ErrNoCustomers = errors.New("failed to get Customers")

type lambdaHandler struct {
	customString string
	logger       *log.Logger
}

// New creates a new handler for Lambda one.
func New(logger *log.Logger, customString string) lambda.Handler {
	return util.NewHandlerV1(lambdaHandler{
		customString: customString,
		logger:       logger,
	})
}

// Handle implements util.LambdaHTTPV1 interface. It contains the logic for the handler.
func (handler lambdaHandler) Handle(ctx context.Context, request *events.APIGatewayProxyRequest) (response *events.APIGatewayProxyResponse, err error) {
	response = &events.APIGatewayProxyResponse{}

	key := request.QueryStringParameters["appkey"]

	app_key := os.Getenv("APP_KEY")

	if key != app_key {
		response.StatusCode = http.StatusBadGateway
		response.Body = string("Invalid APP Key!")

		return response, nil
	}

	id, ok := request.PathParameters["id"]
	if !ok {
		handler.logger.Println("No id given.")
		response.StatusCode = http.StatusBadGateway
	}

	err = deleteCustomer(handler, id)
	if err != nil {
		response.StatusCode = http.StatusBadGateway
		response.Body = string("Error creating or updating post")
	}

	data, err := json.MarshalIndent("", "", "    ")
	if err != nil {
		handler.logger.Print("Failed to JSON marshal response.\nError: %w", err)
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	response.StatusCode = http.StatusAccepted
	response.Body = string(data)

	return response, nil
}

func deleteCustomer(handler lambdaHandler, id string) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		TableName: aws.String("kreativroni.customers"),
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		log.Println(id)
		log.Fatalf("Got error calling DeleteItem: %s", err)
		return err
	}

	return nil
}
