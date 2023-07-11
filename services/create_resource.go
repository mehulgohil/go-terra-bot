package services

import (
	"encoding/json"
	"fmt"
	"log"
)

func CreateResource(userPrompt string) {
	openAiResponse, err := GetFunctionArgumentsFromOpenAI(userPrompt)
	if err != nil {
		aiSummary, err := SummarizeResponseFromOpenAI(userPrompt, openAiResponse.Choices[0].Message.FunctionCall, err.Error())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(aiSummary)
		return
	}

	var jsonMap FunctionCallResponseArguments
	json.Unmarshal([]byte(openAiResponse.Choices[0].Message.FunctionCall.Arguments), &jsonMap)

	azRG := AzureResourceGroup{
		ResourceGroupName:     "rg-demo-001",
		ResourceGroupLocation: "West Europe",
	}
	azRGMarshed, _ := json.Marshal(azRG)

	var outputMessage string
	err = CreateCloudResource(jsonMap.TerraformCloudType, jsonMap.TerraformResourceName, azRG)
	if err != nil {
		outputMessage = "error creating resource. Reason - " + err.Error()
	} else {
		outputMessage = "successfully created resource with details - " + string(azRGMarshed)
	}

	aiSummary, err := SummarizeResponseFromOpenAI(userPrompt, openAiResponse.Choices[0].Message.FunctionCall, outputMessage)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(aiSummary)
}
