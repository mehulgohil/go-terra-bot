package services

import (
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
)

func CreateResource(userPrompt string, dryRun bool) {
	//var ResourceFields = make(map[string][]string)
	//ResourceFields["azurerm_resource_group"] = []string{"name", "location"}
	//ResourceFields["aws_vpc"] = []string{"cidr"}

	openAiResponse, err := GetFunctionArgumentsFromOpenAI(userPrompt)
	if err != nil {
		handleCLIResponse(userPrompt, nil, "error calling openapi function call: "+err.Error())
	}

	var jsonMap FunctionCallResponseArguments
	if len(openAiResponse.Choices) != 0 && openAiResponse.Choices[0].Message.FunctionCall != nil {
		err := json.Unmarshal([]byte(openAiResponse.Choices[0].Message.FunctionCall.Arguments), &jsonMap)
		if err != nil {
			handleCLIResponse(userPrompt, openAiResponse.Choices[0].Message.FunctionCall, "error unmarshalling the openapi function response: "+err.Error())
		}
	} else {
		handleCLIResponse(userPrompt, nil, "please provide a valid prompt. For example `create an aws vpc`")
	}

	azRG := map[string]string{
		"ResourceGroupName":     "rg-demo-001",
		"ResourceGroupLocation": "West Europe",
		"AWS_VPC_CIDR_BLOCK":    "10.0.0.0/16",
	}
	azRGMarshed, _ := json.Marshal(azRG)

	var outputMessage string
	err = CreateCloudResource(jsonMap.TerraformCloudType, jsonMap.TerraformResourceName, azRG, dryRun)
	if err != nil {
		outputMessage = "error creating resource. Reason - " + err.Error()
	} else {
		outputMessage = "successfully created resource with details - " + string(azRGMarshed)
	}

	handleCLIResponse(userPrompt, openAiResponse.Choices[0].Message.FunctionCall, outputMessage)
}

func handleCLIResponse(userPrompt string, funcCall *openai.FunctionCall, output string) {
	aiSummary, err := SummarizeResponseFromOpenAI(userPrompt, funcCall, output)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(aiSummary)
	return
}
