package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

func doMake(arg2, arg3 string) error {
	switch arg2 {
	case "migration":
		dbType := at.DB.Type
		if arg3 == "" {
			exitGracefully(errors.New("you must give the migraion a name"))
		}

		fileName := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), arg3)
		upFile := at.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
		downFile := at.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

		err := copyFileFromTemplate("templates/migrations/"+dbType+".up.sql", upFile)
		if err != nil {
			exitGracefully(err)
		}

		err = copyFileFromTemplate("templates/migrations/"+dbType+".down.sql", downFile)
		if err != nil {
			exitGracefully(err)
		}
	case "auth":
		err := doAuth()
		if err != nil {
			exitGracefully(err)
		}
	case "handler":
		if arg3 == "" {
			exitGracefully(errors.New("you must give the handler a name"))
		}

		fileName := at.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
		if fileExists(fileName) {
			exitGracefully(errors.New(fileName + "already exists!"))
		}

		err := copyFileContentFromTemplate("templates/handler.go.txt", fileName, arg3)
		if err != nil {
			exitGracefully(err)
		}
	default:
		exitGracefully(errors.New("unsupported make arguments"))
	}
	return nil
}
