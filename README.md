# go-terra-bot

`go-terra-bot` is a command-line tool that allows you to create basic cloud resources using OpenAI and Terraform with just one line. It provides a simple and efficient way to interact with cloud services. It also provides the generated resource tf files.

## Installation

To install `go-terra-bot`, follow these steps:

1. Make sure you have Go installed on your system. If not, you can download and install it from the official Go website: [https://golang.org/dl/](https://golang.org/dl/)

2. Open a terminal and run the following command to download and install the tool:

   ```bash
   go install github.com/mehulgohil/go-terra-bot

3. Set OpenAI Key Env Variable
- OPENAPI_KEY: Your OpenAI Token
## Usage

The basic usage of `go-terra-bot` is as follows:

```bash
  go-terra-bot -p "<prompt>"
```

Replace <prompt> with the desired prompt for creating a cloud resource.

```bash
go-terra-bot -p "create an azure resource group"
```

## Cloud Provider Configuration
### AWS
If you want to create "AWS" resources, the following environment variables are required:

- AWS_ACCESS_KEY_ID: AWS access key ID.
- AWS_SECRET_ACCESS_KEY: AWS secret access key.
- AWS_REGION: AWS region.

### Azure
If you want to create "Azure" resources, the following environment variables are required:

- ARM_CLIENT_ID: Azure AD application client ID.
- ARM_CLIENT_SECRET: Azure AD application client secret.
- ARM_TENANT_ID: Azure AD tenant ID.
- ARM_SUBSCRIPTION_ID: Azure subscription ID.

## Current Supported Resources
### AWS
- aws_vpc

### Azure
- azurerm_resource_group