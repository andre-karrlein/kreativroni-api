package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/andre-karrlein/kreativroni-api/model"
	"github.com/andre-karrlein/kreativroni-api/util"
)

// ErrNoProduct indicates that API failed in some way
var ErrNoProductVariation = errors.New("failed to get Product Variation")

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

	id := request.QueryStringParameters["id"]

	if id == "" {
		data, err := json.MarshalIndent(loadAllVariations(language), "", "    ")
		if err != nil {
			handler.logger.Print("Failed to JSON marshal response.\nError: %w", err)
			response.StatusCode = http.StatusInternalServerError
			return response, nil
		}

		response.StatusCode = http.StatusOK
		response.Body = string(data)

		return response, nil
	}

	data, err := json.MarshalIndent(loadVariations(language, id), "", "    ")
	if err != nil {
		handler.logger.Print("Failed to JSON marshal response.\nError: %w", err)
		response.StatusCode = http.StatusInternalServerError
		return response, nil
	}

	response.StatusCode = http.StatusOK
	response.Body = string(data)

	return response, nil
}

func loadVariations(language string, id string) []model.Variation {
	b, err := util.Etsy_request("https://openapi.etsy.com/v3/application/shops/31340310/listings/" + id + "/variation-images")
	if err != nil {
		log.Println(err)
		return nil
	}
	sb := string(b)

	var variationData model.VariationData
	json.Unmarshal([]byte(sb), &variationData)

	var variations []model.Variation
	index := 1

	for _, listingVariation := range variationData.Results {
		variations = append(variations, model.Variation{
			Id:         index,
			PropertyId: listingVariation.Id,
			ValueId:    listingVariation.ValueId,
			Value:      listingVariation.Value,
			ImageId:    listingVariation.ImageId,
		})

		index++
	}

	return variations
}

func loadAllVariations(language string) []model.Variations {
	b, err := util.Etsy_request("https://openapi.etsy.com/v3/application/shops/31340310/listings/active")
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

	var variations []model.Variations
	for _, id := range ids {
		b, err := util.Etsy_request("https://openapi.etsy.com/v3/application/shops/31340310/listings/" + id + "/variation-images")
		if err != nil {
			log.Println(err)
			return nil
		}
		sb := string(b)

		var variationData model.VariationData
		json.Unmarshal([]byte(sb), &variationData)

		var variationElems []model.Variation
		index := 1

		for _, listingVariation := range variationData.Results {
			variationElems = append(variationElems, model.Variation{
				Id:         index,
				PropertyId: listingVariation.Id,
				ValueId:    listingVariation.ValueId,
				Value:      listingVariation.Value,
				ImageId:    listingVariation.ImageId,
			})

			index++
		}

		int_id, _ := strconv.Atoi(id)
		variations = append(variations, model.Variations{
			Id:         int_id,
			Variations: variationElems,
		})
	}

	return variations
}
