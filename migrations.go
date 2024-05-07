package atlas

import (
	"log"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (a *Atlas) MigrateUp(dsn string) error {
	m, err := migrate.New("file://"+a.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		log.Println("Error running migration: ", err)
		return err
	}

	return nil
}

func (a *Atlas) MigrateDownAll(dsn string) error {
	m, err := migrate.New("file://"+a.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Down(); err != nil {
		log.Println("Error rolling migrations back: ", err)
		return err
	}

	return nil
}

func (a *Atlas) Steps(n int, dsn string) error {
	m, err := migrate.New("file://"+a.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Steps(n); err != nil {
		log.Println("Error running migrations steps: ", err)
		return err
	}

	return nil
}

func (a *Atlas) MigrateForce(dsn string) error {
	m, err := migrate.New("file://"+a.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err = m.Force(-1); err != nil {
		log.Println("Error running migrations force: ", err)
		return err
	}

	return nil
}
