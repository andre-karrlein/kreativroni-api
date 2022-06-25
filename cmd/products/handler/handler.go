package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/andre-karrlein/kreativroni-api/model"
	"github.com/andre-karrlein/kreativroni-api/util"
)

// ErrNoProduct indicates that API failed in some way
var ErrNoProduct = errors.New("failed to get Product")

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
	language := request.QueryStringParameters["lang"]

	if language == "" {
		response.StatusCode = http.StatusBadGateway
		response.Body = string("Missing Language!")

		return response, nil
	}

	var products []model.Product
	new := request.QueryStringParameters["new"]
	if new == "yes" {
		products = getProducts(handler, language)
	} else {
		products = loadProducts(language)
	}

	data, err := json.MarshalIndent(products, "", "    ")
	if err != nil {
		handler.logger.Print("Failed to JSON marshal response.\nError: %w", err)
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	response.StatusCode = http.StatusOK
	response.Body = string(data)
	response.Headers = map[string]string{"Access-Control-Allow-Origin": "*"}

	return response, nil
}

func getProducts(handler lambdaHandler, language string) []model.Product {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	out, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("kreativroni"),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				S: aws.String("products"),
			},
			"SK": {
				S: aws.String(language),
			},
		},
	})

	if err != nil {
		panic(err)
	}

	products := model.AwsProduct{}

	err = dynamodbattribute.UnmarshalMap(out.Item, &products)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return products.Products
}

func loadProducts(language string) []model.Product {
	b, err := util.Etsy_request("https://openapi.etsy.com/v3/application/shops/31340310/listings/active?limit=100")
	if err != nil {
		log.Println(err)
		return nil
	}
	sb := string(b)

	var listings model.ListingData
	json.Unmarshal([]byte(sb), &listings)

	var ids []string

	for _, listing := range listings.Results {
		ids = append(ids, strconv.Itoa(listing.Id))
	}

	id_string := strings.Join(ids[:], ",")

	url := "https://openapi.etsy.com/v3/application/listings/batch?limit=100&includes=images&language=" + language + "&listing_ids=" + id_string
	b, err = util.Etsy_request(url)
	if err != nil {
		log.Println(err)
		return nil
	}
	sb = string(b)

	var etsyData model.EtsyProductData
	json.Unmarshal([]byte(sb), &etsyData)

	var products []model.Product
	for _, listingProduct := range etsyData.Results {
		products = append(products, model.Product(listingProduct))
	}

	return products
}
