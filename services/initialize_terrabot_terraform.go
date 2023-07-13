package services

import (
	"log"
	"os"
	"strings"
)

func InitializeTerraformFolders() {
	err := os.Mkdir("terrabot-tf", os.ModePerm)
	if err != nil && !strings.Contains(err.Error(), "Cannot create a file when that file already exists") {
		log.Fatal(err)
	}
}
