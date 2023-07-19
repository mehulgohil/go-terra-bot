package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	goTerraBotPackageName = "github.com/mehulgohil/go-terra-bot"
	Azure                 = "Azure"
	AWS                   = "AWS"
)

var supportedResource = make(map[string][]string)

func init() {
	supportedResource[Azure] = []string{"azurerm_resource_group"}
	supportedResource[AWS] = []string{"aws_vpc"}
}

var supportedProvider = []string{Azure, AWS}

func CreateCloudResource(cloudProvider string, tfResourceName string, resourceParams interface{}, dryRun bool) error {
	fmt.Println(fmt.Sprintf("Trying to create `%s` in %s...", tfResourceName, cloudProvider))

	pkgDir, err := getPackageDirPath(goTerraBotPackageName)
	if err != nil {
		return fmt.Errorf("couldn't find the `go-terra-bot` pkg installed in local: %w", err)
	}

	slashedPkgDir := filepath.ToSlash(pkgDir)
	err = createAndValidateTFFiles(cloudProvider, slashedPkgDir, dryRun)
	if err != nil {
		return errors.New("error creating terraform folder structure: " + err.Error())
	}

	// check supported cloud provider
	if !contains(supportedProvider, cloudProvider) {
		return fmt.Errorf("unsupported cloud provider. Supported providers are %s", strings.Join(supportedProvider, ","))
	}

	// check supported resources
	if !contains(supportedResource[cloudProvider], tfResourceName) {
		return fmt.Errorf("unsupported cloud resource. Supported resources are %s", strings.Join(supportedResource[cloudProvider], ","))
	}

	tmpFile := fmt.Sprintf("%s/templates/%s/%s.tmpl", slashedPkgDir, cloudProvider, tfResourceName)
	tmpl, err := template.ParseFiles(tmpFile)
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("terrabot-tf/%s/%s.tf", cloudProvider, tfResourceName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.Execute(file, resourceParams)
	if err != nil {
		return err
	}

	// if dry run, we exit after creating tf files
	if dryRun {
		return nil
	}

	tf, err := installAndConfigureTerraform(cloudProvider)
	if err != nil {
		return err
	}

	return applyTerraformChanges(tf)
}

func installAndConfigureTerraform(cloudProvider string) (*tfexec.Terraform, error) {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.0.6")),
	}
	execPath, err := installer.Install(context.Background())
	if err != nil {
		return nil, fmt.Errorf("couldn't install terraform: %w", err)
	}

	terraformWorkingDir := fmt.Sprintf("terrabot-tf/%s", cloudProvider)
	tf, err := tfexec.NewTerraform(terraformWorkingDir, execPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialize terraform object: %w", err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		return nil, fmt.Errorf("couldn't run terraform init: %w", err)
	}

	return tf, nil
}

func applyTerraformChanges(tf *tfexec.Terraform) error {
	plan, err := tf.Plan(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't run terraform plan: %w", err)
	}

	if plan {
		err = tf.Apply(context.Background())
		if err != nil {
			return fmt.Errorf("couldn't run terraform apply: %w", err)
		}
	}

	return nil
}

func createAndValidateTFFiles(cloudProvider string, pkgDir string, dryRun bool) error {

	// if not a dry run, we will check all env vars
	if !dryRun {
		switch cloudProvider {
		case AWS:
			awsEnvVars := []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_REGION"}
			for _, envVar := range awsEnvVars {
				if _, ok := os.LookupEnv(envVar); !ok {
					return fmt.Errorf("missing %s environment variable for %s", envVar, AWS)
				}
			}
		case Azure:
			azureEnvVars := []string{"ARM_CLIENT_ID", "ARM_CLIENT_SECRET", "ARM_TENANT_ID", "ARM_SUBSCRIPTION_ID"}
			for _, envVar := range azureEnvVars {
				if _, ok := os.LookupEnv(envVar); !ok {
					return fmt.Errorf("missing %s environment variable for %s", envVar, Azure)
				}
			}
		}
	}

	terrabotDir := filepath.Join("terrabot-tf", cloudProvider)
	err := os.Mkdir(terrabotDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	filePath := filepath.Join(terrabotDir, "auth.tf")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpFile := filepath.Join(pkgDir, "templates", cloudProvider, "auth.tmpl")
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
