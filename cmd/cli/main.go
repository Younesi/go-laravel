package main

import (
	"errors"
	"os"

	"github.com/fatih/color"
	"github.com/younesi/atlas"
)

var at atlas.Atlas

const version = "0.0.1"

func main() {
	var message string
	arg1, arg2, arg3, err := validateInput()

	if err != nil {
		exitGracefully(err)
	}

	setup()

	switch arg1 {
	case "make":
		if arg2 == "" {
			exitGracefully(errors.New("makes requies a subcommand: (migrations|model|handler)"))
		}

		err = doMake(arg2, arg3)
		if err != nil {
			exitGracefully(err)
		}
	case "migrate":
		if arg2 == "" {
			arg2 = "up"
		}
		err := doMigrate(arg2, arg3)
		if err != nil {
			exitGracefully(err)
		}
		message = "Migrations complete!"
	case "help":
		showHelp()
	case "version":
		color.Yellow("Application version:" + version)
	default:
		showHelp()
	}

	exitGracefully(nil, message)
}

func validateInput() (string, string, string, error) {
	var arg1, arg2, arg3 string

	if len(os.Args) > 1 {
		arg1 = os.Args[1]

		if len(os.Args) >= 3 {
			arg2 = os.Args[2]
		}
		if len(os.Args) >= 4 {
			arg3 = os.Args[3]
		}
	} else {
		color.Red("Error: command required")
		showHelp()
		return "", "", "", errors.New("command required")
	}

	return arg1, arg2, arg3, nil
}

func exitGracefully(err error, msg ...string) {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}

	if err != nil {
		color.Red("Error : %v\n", err)
	}

	if len(message) > 0 {
		color.Yellow(message)
	} else {
		color.Green("finished")
	}

	os.Exit(0)
}
