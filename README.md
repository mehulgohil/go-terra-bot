# go-terra-bot

`go-terra-bot` is a command-line tool that allows you to create basic cloud resources using OpenAI and Terraform with just one line. It provides a simple and efficient way to interact with cloud services.

## Installation

To install `go-terra-bot`, follow these steps:

1. Make sure you have Go installed on your system. If not, you can download and install it from the official Go website: [https://golang.org/dl/](https://golang.org/dl/)

2. Open a terminal and run the following command to download and install the tool:

   ```bash
   go get -u github.com/mehulgohil/go-terra-bot/cmd/go-terra-bot

## Usage

The basic usage of `go-terra-bot` is as follows:

```bash
  go-terra-bot -p "<prompt>"
```

Replace <prompt> with the desired prompt for creating a cloud resource.

```bash
go-terra-bot -p "create an azure resource group"
```