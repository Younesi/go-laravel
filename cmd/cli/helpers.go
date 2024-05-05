package main

import (
	"os"

	"github.com/joho/godotenv"
)

func setup() {
	err := godotenv.Load()
	if err != nil {
		exitGracefully(err)
	}
	path, err := os.Getwd()
	if err != nil {
		exitGracefully(err)
	}
	at.RootPath = path
	at.DB.Type = os.Getenv("DATABASE_TYPE")
}
