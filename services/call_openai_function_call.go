package services

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"os"
	"sync"
)

var client *openai.Client
var runOnce sync.Once

func CallOpenAIFunctionCall(userPrompt string) {
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
					Name:        "create_azure_resource",
					Description: "Create an azure resource",
					Parameters: RequestParameters{
						Type: string(jsonschema.Object),
						Properties: RequestParameterProperty{
							TerraformResourceName: TFResourceName{
								Type:        "string",
								Description: "The terraform resource name, e.g. azurerm_resource_group",
							},
						},
						Required: []string{"terraform_resource_name"},
					},
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.FunctionCall)
}

func getClient() {
	runOnce.Do(func() {
		client = openai.NewClient(os.Getenv("OPENAPI_KEY"))
	})
}
