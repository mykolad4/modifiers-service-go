package main

import (
	"context"
	"encoding/json"
	"github.com/OrderAhead/FastaServerless/services/modifier-api/handler/schema"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/graphql-go/graphql"
)

type RequestBody struct {
	Query          string                 `json:"query"`
	VariableValues map[string]interface{} `json:"variables"`
	OperationName  string                 `json:"operationName"`
}

func executeQuery(request RequestBody, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:         schema,
		VariableValues: request.VariableValues,
		RequestString:  request.Query,
		OperationName:  request.OperationName,
	})

	return result
}

// Handler of HTTP event
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	requestBody := RequestBody{}
	err := json.Unmarshal([]byte(request.Body), &requestBody)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, err
	}

	graphQLResult := executeQuery(requestBody, schema.Schema)
	responseJSON, err := json.Marshal(graphQLResult)

	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, err
	}

	return events.APIGatewayProxyResponse{Body: string(responseJSON[:]), StatusCode: 200}, nil
}

//type FlatModifier struct {
//	merchant string
//	group string
//	product string
//	gAtLeast int
//	gAtMost int
//	name string
//	cost float64
//	atLeast int
//	atMost int
//	isDefault bool
//	isHidden bool
//}

//type FlatModifier struct {
//	Merchant string `dynamodbav:"merchant"`
//	Group string `dynamodbav:"group"`
//	Product string `dynamodbav:"product"`
//	GAtLeast int `dynamodbav:"g_at_least"`
//	GAtMost int `dynamodbav:"g_at_most"`
//	Name string `dynamodbav:"name"`
//	Cost float64 `dynamodbav:"cost"`
//	AtLeast int `dynamodbav:"at_least"`
//	AtMost int `dynamodbav:"at_most"`
//	IsDefault bool  `dynamodbav:"is_default"`
//	IsHidden bool  `dynamodbav:"is_hidden"`
//}

func main() {
	lambda.Start(Handler)
	//flatModifier := FlatModifier {
	//	Merchant: "merchant",
	//	Product: "product",
	//	Group: "group",
	//	GAtLeast: 1,
	//	GAtMost: 2,
	//	Name: "name",
	//	Cost: 77.77,
	//	AtLeast: 3,
	//	AtMost: 4,
	//	IsDefault: true,
	//	IsHidden: false,
	//}
	//
	//av, err := attributevalue.MarshalMap(flatModifier)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(av)
}
