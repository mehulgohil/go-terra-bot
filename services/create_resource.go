package services

import (
	"encoding/json"
	"fmt"
	"log"
)

func CreateResource(userPrompt string, dryRun bool) {
	openAiResponse, err := GetFunctionArgumentsFromOpenAI(userPrompt)
	if err != nil {
		fmt.Println(err)
		aiSummary, err := SummarizeResponseFromOpenAI(userPrompt, openAiResponse.Choices[0].Message.FunctionCall, err.Error())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(aiSummary)
		return
	}

	var jsonMap FunctionCallResponseArguments
	json.Unmarshal([]byte(openAiResponse.Choices[0].Message.FunctionCall.Arguments), &jsonMap)

	azRG := map[string]string{
		"ResourceGroupName":     "rg-demo-001",
		"ResourceGroupLocation": "West Europe",
		"AWS_VPC_CIDR_BLOCK":    "10.0.0.0/16",
	}
	azRGMarshed, _ := json.Marshal(azRG)

	var outputMessage string
	err = CreateCloudResource(jsonMap.TerraformCloudType, jsonMap.TerraformResourceName, azRG, dryRun)
	if err != nil {
		fmt.Println(err)
		outputMessage = "error creating resource. Reason - " + err.Error()
	} else {
		outputMessage = "successfully created resource with details - " + string(azRGMarshed)
	}

	aiSummary, err := SummarizeResponseFromOpenAI(userPrompt, openAiResponse.Choices[0].Message.FunctionCall, outputMessage)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	fmt.Println(aiSummary)
}
