package schema

import (
	"fmt"
	"github.com/OrderAhead/FastaServerless/services/modifier-api/handler/services"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/graphql-go/graphql"
	"log"
	"strconv"
)

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"getModifiersByMerchant": &graphql.Field{
			Type: &graphql.List{OfType: graphQLModifierGroupType},
			Description: "Get modifiers by merchant",
			Args: graphql.FieldConfigArgument{
				"merchant": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				merchant, _ := params.Args["merchant"].(string)
				result, err := services.GetModifiersByMerchant(merchant)

				if err != nil {
					log.Println("error:" + err.Error())
					return nil, err
				}

				response := make(map[string]ModifierGroup)

				for _, item := range result.Items {
					cost,_ := strconv.ParseFloat(item["cost"].(*types.AttributeValueMemberN).Value, 2)
					atLeast,_ := strconv.Atoi(item["atLeast"].(*types.AttributeValueMemberN).Value)
					atMost,_ := strconv.Atoi(item["atMost"].(*types.AttributeValueMemberN).Value)
					name := item["name"].(*types.AttributeValueMemberS).Value
					isDefault := item["isDefault"].(*types.AttributeValueMemberBOOL).Value
					isHidden := item["isHidden"].(*types.AttributeValueMemberBOOL).Value
					group := item["group"].(*types.AttributeValueMemberS).Value

					modifier := Modifier{
						Name: name,
						Cost: cost,
						AtLeast: atLeast,
						AtMost: atMost,
						IsDefault: isDefault,
						IsHidden: isHidden,
					}

					if val, ok := response[group]; ok {
						val.Modifiers = append(val.Modifiers, &modifier)
						response[group] = val
					} else {
						merchant := item["merchant"].(*types.AttributeValueMemberS).Value
						product := item["product"].(*types.AttributeValueMemberS).Value
						gAtLeast,_ := strconv.Atoi(item["gAtLeast"].(*types.AttributeValueMemberN).Value)
						gAtMost,_ := strconv.Atoi(item["gAtMost"].(*types.AttributeValueMemberN).Value)

						groupModifier := ModifierGroup{
							Group: group,
							Merchant: merchant,
							Product: product,
							AtLeast: gAtLeast,
							AtMost: gAtMost,
						}

						groupModifier.Modifiers = append(groupModifier.Modifiers, &modifier)
						response[group] = groupModifier
					}
				}

				values := make([]ModifierGroup, 0, len(response))

				for _, v := range response {
					values = append(values, v)
				}

				return values, err
			},
		},
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"createModifier": &graphql.Field{
			Type: graphQLModifierGroupType,
			Args: graphql.FieldConfigArgument{
				"merchant": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"product": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"group": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"gAtLeast": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"gAtMost": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"cost": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
				"atLeast": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"atMost": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"isDefault": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
				"isHidden": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				flatModifier := FlatModifier{
					Merchant: params.Args["merchant"].(string),
					Product: params.Args["product"].(string),
					Group: params.Args["group"].(string),
					GAtLeast: params.Args["gAtLeast"].(int),
					GAtMost: params.Args["gAtMost"].(int),
					Name: params.Args["name"].(string),
					Cost: params.Args["cost"].(float64),
					AtLeast: params.Args["atLeast"].(int),
					AtMost: params.Args["atMost"].(int),
					IsDefault: params.Args["isDefault"].(bool),
					IsHidden: params.Args["isHidden"].(bool),
				}

				log.Println(flatModifier)

				av, err := attributevalue.MarshalMap(flatModifier)
				if err != nil {
					log.Fatalf("Got error marshalling new modifier item: %s", err)
				}

				log.Println(av)

				res, err := services.AddModifier(av)
				if res == nil || err != nil {
					fmt.Println("Error==========" + err.Error())
				}

				modifier := Modifier{
					Name: params.Args["name"].(string),
					Cost: params.Args["cost"].(float64),
					AtLeast: params.Args["atLeast"].(int),
					AtMost: params.Args["atMost"].(int),
					IsDefault: params.Args["isDefault"].(bool),
					IsHidden: params.Args["isHidden"].(bool),
				}
				group := ModifierGroup{
					Merchant: params.Args["merchant"].(string),
					Product: params.Args["product"].(string),
					Group: params.Args["group"].(string),
					AtLeast: params.Args["gAtLeast"].(int),
					AtMost: params.Args["gAtMost"].(int),
					Modifiers: []*Modifier{&modifier},
				}

				return group, nil
			},
		},
	},
})

// Schema - GraphQL root schema
var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
