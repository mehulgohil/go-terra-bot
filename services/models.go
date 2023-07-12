package services

type FunctionCallResponseArguments struct {
	TerraformResourceName string `json:"terraform_resource_name"`
	TerraformCloudType    string `json:"terraform_cloud_Provider"`
}

type AzureResourceGroup struct {
	ResourceGroupName     string
	ResourceGroupLocation string
}
