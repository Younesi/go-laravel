package main

func doMigrate(arg2, arg3 string) error {
	dsn := getDSN()

	switch arg2 {
	case "up":
		err := at.MigrateUp(dsn)
		if err != nil {
			return err
		}
	case "down":
		if arg3 == "all" {
			err := at.MigrateDownAll(dsn)
			if err != nil {
				return err
			}
		}

		err := at.Steps(-1, dsn)
		if err != nil {
			return err
		}
	case "reset":
		err := at.MigrateDownAll(dsn)
		if err != nil {
			return err
		}
		err = at.MigrateUp(dsn)
		if err != nil {
			return err
		}
	default:
		showHelp()
	}

	return nil
}
