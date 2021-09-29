package schema

import "github.com/graphql-go/graphql"

type Modifier struct {
	Name string `json:"name"`
	Cost float64 `json:"cost"`
	AtLeast int `json:"at_least"`
	AtMost int `json:"at_most"`
	IsDefault bool `json:"is_default"`
	IsHidden bool `json:"is_hidden"`
}

type ModifierGroup struct {
	Merchant string `json:"merchant"`
	Group string `json:"group"`
	Product string `json:"product"`
	AtLeast int `json:"at_least"`
	AtMost int `json:"at_most"`
	Modifiers []*Modifier `json:"modifiers"`
}

type FlatModifier struct {
	Merchant string `dynamodbav:"merchant"`
	Group string `dynamodbav:"group"`
	Product string `dynamodbav:"product"`
	GAtLeast int `dynamodbav:"gAtLeast"`
	GAtMost int `dynamodbav:"gAtMost"`
	Name string `dynamodbav:"name"`
	Cost float64 `dynamodbav:"cost"`
	AtLeast int `dynamodbav:"atLeast"`
	AtMost int `dynamodbav:"atMost"`
	IsDefault bool  `dynamodbav:"isDefault"`
	IsHidden bool  `dynamodbav:"isHidden"`
}

var graphQLModifierType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Modifier",
	Description: "This represents modifier",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"cost": &graphql.Field{
			Type: graphql.Float,
		},
		"atLeast": &graphql.Field{
			Type: graphql.Int,
		},
		"atMost": &graphql.Field{
			Type: graphql.Int,
		},
		"isDefault": &graphql.Field{
			Type: graphql.Boolean,
		},
		"isHidden": &graphql.Field{
			Type: graphql.Boolean,
		},
	},
})

var graphQLModifierGroupType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ModifierGroup",
	Description: "This represents group of modifiers",
	Fields: graphql.Fields{
		"merchant": &graphql.Field{
			Type: graphql.String,
		},
		"product": &graphql.Field{
			Type: graphql.String,
		},
		"group": &graphql.Field{
			Type: graphql.String,
		},
		"atLeast": &graphql.Field{
			Type: graphql.Int,
		},
		"atMost": &graphql.Field{
			Type: graphql.Int,
		},
		"modifiers": &graphql.Field{
			Type: &graphql.List{OfType: graphQLModifierType},
		},
	},
})