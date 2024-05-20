package atlas

import "database/sql"

type Database struct {
	Type string
	Pool *sql.DB
}

type initPath struct {
	rootPath    string
	folderNames []string
}

type cookieConfig struct {
	name     string
	lifetime string
	persist  string // does it persist between browser's closes
	secure   string
	domain   string
}

type databaseConfig struct {
	dsn      string
	database string
}

type redisConfig struct {
	host     string
	password string
	prefix   string
}
