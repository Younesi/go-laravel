package main

import (
	"embed"
	"errors"
	"io/ioutil"
	"os"
)

//go:embed templates
var templateFS embed.FS

func copyFileFromTemplate(templatePath, targetFile string) error {
	if fileExists(targetFile) {
		exitGracefully(errors.New(targetFile + "already exists!"))
	}

	data, err := templateFS.ReadFile(templatePath)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile(data, targetFile)
	if err != nil {
		exitGracefully(err)
	}

	return nil
}

func copyDataToFile(data []byte, to string) error {
	err := ioutil.WriteFile(to, data, 0644)
	if err != nil {
		exitGracefully(err)
	}

	return nil
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
