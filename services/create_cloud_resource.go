package services

import (
	"errors"
	"os"
	"strings"
	"text/template"
)

var supportedResource = []string{"azurerm_resource_group"}

func CreateCloudResource(cloudProvider string, tfResourceName string, resourceParams interface{}) error {
	if !contains(supportedResource, tfResourceName) {
		return errors.New("Unsupported cloud resource. Supported resources are " + strings.Join(supportedResource, ","))
	}

	var tmpFile = "./templates/" + cloudProvider + "/" + tfResourceName + ".tmpl"

	tmpl, err := template.ParseFiles(tmpFile)
	if err != nil {
		return err
	}

	file, err := os.Create("output.tf")
	if err != nil {
		os.Exit(1)
	}
	defer file.Close()

	err = tmpl.Execute(file, resourceParams)
	if err != nil {
		return err
	}

	//terraform init
	//terraform plan
	//terraform apply

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
