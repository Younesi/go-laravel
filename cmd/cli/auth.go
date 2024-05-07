package main

import (
	"fmt"
	"time"
)

func doAuth() error {
	dbType := at.DB.Type
	fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())
	upFile := at.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
	downFile := at.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

	err := copyFileFromTemplate("templates/migrations/auth_tables."+dbType+".up.sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	err = copyFileFromTemplate("templates/migrations/auth_tables."+dbType+".down.sql", downFile)
	if err != nil {
		exitGracefully(err)
	}
	return nil
}
