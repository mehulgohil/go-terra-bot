package services

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	goTerraBotPackageName = "github.com/mehulgohil/go-terra-bot"
)

var supportedResource = []string{"azurerm_resource_group", "aws_vpc"}
var supportedProvider = []string{"Azure", "AWS"}

func CreateCloudResource(cloudProvider string, tfResourceName string, resourceParams interface{}, dryRun bool) error {
	pkgDir, err := getPackageDirPath(goTerraBotPackageName)
	if err != nil {
		return errors.New("couldn't find the `go-terra-bot` pkg installed in local: " + err.Error())
	}
	slashedPkgDir := filepath.ToSlash(pkgDir)

	err = createAndValidateTFFiles(cloudProvider, slashedPkgDir, dryRun)
	if err != nil {
		return errors.New("error creating terraform folder structure: " + err.Error())
	}
	if !contains(supportedProvider, cloudProvider) {
		return errors.New("Unsupported cloud provider. Supported providers are " + strings.Join(supportedResource, ","))
	}
	if !contains(supportedResource, tfResourceName) {
		return errors.New("Unsupported cloud resource. Supported resources are " + strings.Join(supportedResource, ","))
	}

	var tmpFile = slashedPkgDir + "/templates/" + cloudProvider + "/" + tfResourceName + ".tmpl"

	tmpl, err := template.ParseFiles(tmpFile)
	if err != nil {
		return err
	}

	file, err := os.Create("terrabot-tf/" + cloudProvider + "/" + tfResourceName + ".tf")
	if err != nil {
		return err
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

func createAndValidateTFFiles(cloudProvider string, pkgDir string, dryRun bool) error {

	// if not a dry run, we will check all env vars
	if !dryRun {
		switch cloudProvider {
		case "AWS":
			_, ok := os.LookupEnv("AWS_ACCESS_KEY_ID")
			if !ok {
				return errors.New("missing AWS_ACCESS_KEY_ID environment variable for AWS")
			}
			_, ok = os.LookupEnv("AWS_SECRET_ACCESS_KEY")
			if !ok {
				return errors.New("missing AWS_SECRET_ACCESS_KEY environment variable for AWS")
			}
			_, ok = os.LookupEnv("AWS_REGION")
			if !ok {
				return errors.New("missing AWS_REGION environment variable for AWS")
			}
		case "Azure":
			_, ok := os.LookupEnv("ARM_CLIENT_ID")
			if !ok {
				return errors.New("missing ARM_CLIENT_ID environment variable for Azure")
			}
			_, ok = os.LookupEnv("ARM_CLIENT_SECRET")
			if !ok {
				return errors.New("missing ARM_CLIENT_SECRET environment variable for Azure")
			}
			_, ok = os.LookupEnv("ARM_TENANT_ID")
			if !ok {
				return errors.New("missing ARM_TENANT_ID environment variable for Azure")
			}
			_, ok = os.LookupEnv("ARM_SUBSCRIPTION_ID")
			if !ok {
				return errors.New("missing ARM_SUBSCRIPTION_ID environment variable for Azure")
			}
		}
	}

	err := os.Mkdir("terrabot-tf/"+cloudProvider, os.ModePerm)
	if err != nil && !strings.Contains(err.Error(), "Cannot create a file when that file already exists") {
		return err
	}

	file, err := os.Create("terrabot-tf/" + cloudProvider + "/auth.tf")
	if err != nil {
		return err
	}
	defer file.Close()

	var tmpFile = pkgDir + "/templates/" + cloudProvider + "/auth.tmpl"

	tmpl, err := template.ParseFiles(tmpFile)
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, nil)
	if err != nil {
		return err
	}

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

func getPackageDirPath(pkgName string) (string, error) {

	cmd := exec.Command("go", "list", "-f", "{{.Dir}}", pkgName)
	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	dirPath := strings.TrimSpace(string(output))

	return dirPath, nil
}
