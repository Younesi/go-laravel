package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
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

func showHelp() {
	color.Yellow(`Available commands :
		migration <name>		- create a migration file with the specified name
		migrate					- runs all ready-to-run migrations
		migrate	down			- rollbacks the last migration
		migrate	reset			- rollbacks all the migrations
		auth					- create authentication skeletion
		help					- show the help command
		version					- print application version
	`)
}

func getDSN() string {
	dbType := at.DB.Type
	if dbType == "pgx" {
		dbType = "postgres"
	}

	if dbType == "postgres" {
		var dsn string
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASS"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		}

		return dsn
	}

	return "mysql://" + at.BuildDSN()
}
