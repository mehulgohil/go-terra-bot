package services

type RequestParameters struct {
	Type       string                   `json:"type"`
	Properties RequestParameterProperty `json:"properties"`
	Required   []string                 `json:"required"`
}

type RequestParameterProperty struct {
	TerraformResourceName TFResourceName `json:"terraform_resource_name"`
}

type TFResourceName struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}
