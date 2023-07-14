package services

import (
	"log"
	"os"
)

func InitializeTerraformFolders() {
	err := os.Mkdir("terrabot-tf", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}
