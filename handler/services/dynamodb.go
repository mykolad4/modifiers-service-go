package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func createSession() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = "us-east-1"
		return nil
	})
	if err != nil {
		panic(err)
	}

	return dynamodb.NewFromConfig(cfg)
}

func GetModifiersByMerchant(merchant string) (*dynamodb.QueryOutput, error) {
	svc := createSession()

	return svc.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String("modifierTable"),
		KeyConditionExpression: aws.String("merchant = :merchant"),
		ExpressionAttributeValues: map[string] types.AttributeValue{
			":merchant":  &types.AttributeValueMemberS{Value: merchant},
		},
	})
}

func AddModifier(modifier map[string]types.AttributeValue) (*dynamodb.PutItemOutput, error) {
	svc := createSession()

	return svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		Item:      modifier,
		TableName: aws.String("modifierTable"),
	})

}
