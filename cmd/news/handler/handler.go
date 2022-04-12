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

// ErrNoNews indicates that API failed in some way
var ErrNoNews = errors.New("failed to get News")

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
	products_key := os.Getenv("PRODUCTS_KEY")

	if key != app_key && key != products_key {
		response.StatusCode = http.StatusBadGateway
		response.Body = string("Invalid APP Key!")

		return response, nil
	}
	id, ok := request.PathParameters["id"]
	if ok && id != "" {
		data, err := json.MarshalIndent(getSpecificNews(handler, id), "", "    ")
		if err != nil {
			handler.logger.Print("Failed to JSON marshal response.\nError: %w", err)
			response.StatusCode = http.StatusInternalServerError
			return response, nil
		}

		response.StatusCode = http.StatusOK
		response.Body = string(data)

		return response, nil
	}

	all := false

	editor := request.QueryStringParameters["editor"]
	if editor != "" {
		all = true
	}

	data, err := json.MarshalIndent(getAllNews(handler, all), "", "    ")
	if err != nil {
		handler.logger.Print("Failed to JSON marshal response.\nError: %w", err)
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	response.StatusCode = http.StatusOK
	response.Body = string(data)

	return response, nil
}

func getAllNews(handler lambdaHandler, all bool) []model.News {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	expr, err := expression.NewBuilder().WithFilter(
		expression.Or(
			expression.Equal(expression.Name("available"), expression.Value(true)),
			expression.Equal(expression.Name("available"), expression.Value(false)),
		),
	).Build()
	if err != nil {
		panic(err)
	}

	if !all {
		expr, err = expression.NewBuilder().WithFilter(
			expression.Equal(expression.Name("available"), expression.Value(true)),
		).Build()
		if err != nil {
			panic(err)
		}
	}

	out, err := svc.Scan(&dynamodb.ScanInput{
		TableName:                 aws.String("kreativroni.news"),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	if err != nil {
		panic(err)
	}

	news := []model.News{}
	for _, s := range out.Items {
		item := model.News{}

		err = dynamodbattribute.UnmarshalMap(s, &item)
		if err != nil {
			panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
		}
		news = append(news, item)
	}

	return news
}

func getSpecificNews(handler lambdaHandler, id string) model.News {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	out, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("kreativroni.news"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})

	if err != nil {
		panic(err)
	}

	news := model.News{}

	err = dynamodbattribute.UnmarshalMap(out.Item, &news)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return news
}
