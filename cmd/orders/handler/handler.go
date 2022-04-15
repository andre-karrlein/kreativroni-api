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
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"

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
	id, ok := request.PathParameters["id"]
	if ok && id != "" {
		data, err := json.MarshalIndent(getSpecificOrders(handler, id), "", "    ")
		if err != nil {
			handler.logger.Print("Failed to JSON marshal response.\nError: %w", err)
			response.StatusCode = http.StatusInternalServerError
			return response, nil
		}

		response.StatusCode = http.StatusOK
		response.Body = string(data)

		return response, nil
	}

	data, err := json.MarshalIndent(getAllOrders(handler), "", "    ")
	if err != nil {
		handler.logger.Print("Failed to JSON marshal response.\nError: %w", err)
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	response.StatusCode = http.StatusOK
	response.Body = string(data)

	return response, nil
}

func getAllOrders(handler lambdaHandler) []model.Order {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	out, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("kreativroni.orders"),
	})

	if err != nil {
		panic(err)
	}

	orders := []model.Order{}
	for _, s := range out.Items {
		item := model.Order{}

		err = dynamodbattribute.UnmarshalMap(s, &item)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
		orders = append(orders, item)
	}

	return orders
}

func getSpecificOrders(handler lambdaHandler, id string) []model.Order {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	expr, err := expression.NewBuilder().WithFilter(
		expression.Equal(expression.Name("user"), expression.Value(id)),
	).Build()
	if err != nil {
		panic(err)
	}

	out, err := svc.Scan(&dynamodb.ScanInput{
		TableName:                 aws.String("kreativroni.orders"),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	if err != nil {
		panic(err)
	}

	orders := []model.Order{}
	for _, s := range out.Items {
		item := model.Order{}

		err = dynamodbattribute.UnmarshalMap(s, &item)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
		orders = append(orders, item)
	}

	return orders
}
