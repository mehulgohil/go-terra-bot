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
								Description: "The terraform resource name, e.g. azurerm_resource_group, aws_vpc, google_container_aws_cluster",
							},
							"terraform_cloud_Provider": {
								Type: jsonschema.String,
								Enum: []string{"Azure", "AWS"},
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
								Description: "The terraform resource name, e.g. azurerm_resource_group, aws_vpc, google_container_aws_cluster",
							},
							"terraform_cloud_Provider": {
								Type: jsonschema.String,
								Enum: []string{"Azure", "AWS"},
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
