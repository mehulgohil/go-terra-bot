package services

type FunctionCallResponseArguments struct {
	TerraformResourceName string             `json:"terraform_resource_name"`
	TerraformCloudType    string             `json:"terraform_cloud_Provider"`
	ExtraParams           []ExtraParamStruct `json:"extra_params"`
}

type ExtraParamStruct struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
