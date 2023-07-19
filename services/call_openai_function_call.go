package services

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"os"
	"sync"
)

var client *openai.Client
var runOnce sync.Once

func GetFunctionArgumentsFromOpenAI(userPrompt string) (openai.ChatCompletionResponse, error) {
	getClient()

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
			Functions: []openai.FunctionDefinition{
				{
					Name:        "create_cloud_resource",
					Description: "Create a cloud resource",
					Parameters: jsonschema.Definition{
						Type: jsonschema.Object,
						Properties: map[string]jsonschema.Definition{
							"terraform_resource_name": {
								Type:        jsonschema.String,
								Description: "The terraform resource name, e.g. azurerm_resource_group, aws_vpc",
							},
							"terraform_cloud_Provider": {
								Type: jsonschema.String,
								Enum: []string{Azure, AWS},
							},
							"extra_params": {
								Type:        jsonschema.Array,
								Description: "An array of object to pass to function for configuration arguments  of the terraform resource type, e.g. if prompt is `create an azure resource group test in West India`, extra_params should be [{'name': 'name', 'value': 'test'}, {'name':'location', 'value': 'westindia'}, similarly append any other configuration arguments provided by user. The configuration arguments provided should be valid for that terraform resource type.",
								Items: &jsonschema.Definition{
									Type: jsonschema.Object,
									Properties: map[string]jsonschema.Definition{
										"name": {
											Type:        jsonschema.String,
											Description: "terraform configuration argument key field",
										},
										"value": {
											Type:        jsonschema.String,
											Description: "terraform configuration argument value field",
										},
									},
								},
							},
						},
						Required: []string{"terraform_resource_name", "terraform_cloud_provider"},
					},
				},
			},
		},
	)

	if err != nil {
		return openai.ChatCompletionResponse{}, err
	}

	return resp, nil
}

func SummarizeResponseFromOpenAI(userPrompt string, assistantFunctionCall *openai.FunctionCall, functionOutput string) (string, error) {
	getClient()

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
				{
					Role:         openai.ChatMessageRoleAssistant,
					FunctionCall: assistantFunctionCall,
				},
				{
					Role:    openai.ChatMessageRoleFunction,
					Name:    "create_cloud_resource",
					Content: functionOutput,
				},
			},
			Functions: []openai.FunctionDefinition{
				{
					Name:        "create_cloud_resource",
					Description: "Create a cloud resource",
					Parameters: jsonschema.Definition{
						Type: jsonschema.Object,
						Properties: map[string]jsonschema.Definition{
							"terraform_resource_name": {
								Type:        jsonschema.String,
								Description: "The terraform resource name, e.g. azurerm_resource_group, aws_vpc",
							},
							"terraform_cloud_Provider": {
								Type: jsonschema.String,
								Enum: []string{Azure, AWS},
							},
							"extra_params": {
								Type:        jsonschema.Array,
								Description: "An array of dictionary to pass to function for configuration arguments  of the terraform resource type, e.g. if prompt is `create a resource group test in West India`, extra_params should be [{'name': 'name', 'value': 'test'}, {'name':'location', 'value': 'westindia'}, similarly append any other configuration arguments provided by user. The configuration arguments provided should be valid for that terraform resource type.",
								Items: &jsonschema.Definition{
									Type: jsonschema.Object,
									Properties: map[string]jsonschema.Definition{
										"name": {
											Type:        jsonschema.String,
											Description: "terraform configuration argument key field",
										},
										"value": {
											Type:        jsonschema.String,
											Description: "terraform configuration argument value field",
										},
									},
								},
							},
						},
						Required: []string{"terraform_resource_name", "terraform_cloud_provider"},
					},
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func getClient() {
	runOnce.Do(func() {
		client = openai.NewClient(os.Getenv("OPENAPI_KEY"))
	})
}
