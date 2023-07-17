package services

type FunctionCallResponseArguments struct {
	TerraformResourceName string   `json:"terraform_resource_name"`
	TerraformCloudType    string   `json:"terraform_cloud_Provider"`
	ExtraParams           []string `json:"extra_params"`
}
