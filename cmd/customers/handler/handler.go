package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	if ok && id != "" {
		data, err := json.MarshalIndent(getSpecificCustomer(handler, id), "", "    ")
		if err != nil {
			handler.logger.Print("Failed to JSON marshal response.\nError: %w", err)
			response.StatusCode = http.StatusInternalServerError
			return response, nil
		}

		response.StatusCode = http.StatusOK
		response.Body = string(data)

		return response, nil
	}

	data, err := json.MarshalIndent(getAllCustomers(handler), "", "    ")
	if err != nil {
		handler.logger.Print("Failed to JSON marshal response.\nError: %w", err)
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	response.StatusCode = http.StatusOK
	response.Body = string(data)

	return response, nil
}

func getAllCustomers(handler lambdaHandler) []model.Customer {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	out, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("kreativroni.customers"),
	})

	if err != nil {
		panic(err)
	}

	customers := []model.Customer{}
	for _, s := range out.Items {
		item := model.Customer{}

		err = dynamodbattribute.UnmarshalMap(s, &item)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
		customers = append(customers, item)
	}

	return customers
}

func getSpecificCustomer(handler lambdaHandler, id string) model.Customer {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	out, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("kreativroni.customers"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		panic(err)
	}

	customer := model.Customer{}

	err = dynamodbattribute.UnmarshalMap(out.Item, &customer)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return customer
}
