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
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/andre-karrlein/kreativroni-api/model"
	"github.com/andre-karrlein/kreativroni-api/util"
)

// ErrNoOrders indicates that API failed in some way
var ErrNoOrders = errors.New("failed to get Orders")

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

	o := model.Order{}

	json.Unmarshal([]byte(request.Body), &o)

	err = saveOrder(handler, o)
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

	response.StatusCode = http.StatusCreated
	response.Body = string(data)

	return response, nil
}

func saveOrder(handler lambdaHandler, order model.Order) error {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(order)
	if err != nil {
		log.Fatalf("Got error marshalling order item: %s", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("kreativroni.orders"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
