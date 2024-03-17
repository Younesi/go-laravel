package celeritas

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
