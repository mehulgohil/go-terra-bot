package services

type RequestParameters struct {
	Type       string                   `json:"type"`
	Properties RequestParameterProperty `json:"properties"`
	Required   []string                 `json:"required"`
}

type RequestParameterProperty struct {
	TerraformResourceName  TFResourceName `json:"terraform_resource_name"`
	TerraformCloudProvider TFCloudType    `json:"terraform_cloud_Provider"`
}

type TFResourceName struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}
type TFCloudType struct {
	Type string   `json:"type"`
	Enum []string `json:"enum"`
}

type FunctionCallResponseArguments struct {
	TerraformResourceName string `json:"terraform_resource_name"`
	TerraformCloudType    string `json:"terraform_cloud_Provider"`
}

type AzureResourceGroup struct {
	ResourceGroupName     string
	ResourceGroupLocation string
}
