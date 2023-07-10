package main

import (
	"github.com/mehulgohil/go-terra-bot/cmd"
	"os"
)

func main() {
	_, ok := os.LookupEnv("OPENAPI_KEY")
	if !ok {
		panic("missing OPENAPI_KEY environment variable")
	}
	cmd.Execute()
}
