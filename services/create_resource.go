package services

import (
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
)

func CreateResource(userPrompt string, dryRun bool) {
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
	azRG := map[string]string{}
	//add default values
	switch jsonMap.TerraformResourceName {
	case "aws_vpc":
		azRG["cidr_block"] = "10.0.0.0/16"
		for _, eachParam := range jsonMap.ExtraParams {
			if eachParam.Name == "cidr_block" {
				azRG["cidr_block"] = eachParam.Value
			}
		}
	case "azurerm_resource_group":
		azRG["name"] = "rg-demo-001"
		azRG["location"] = "West Europe"
		for _, eachParam := range jsonMap.ExtraParams {
			if eachParam.Name == "name" {
				azRG["name"] = eachParam.Value
			}
			if eachParam.Name == "location" {
				azRG["location"] = eachParam.Value
			}
		}
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
