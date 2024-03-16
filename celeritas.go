package celeritas

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"

	"github.com/go-chi/chi/v5"
	"github.com/younesi/celeritas/render"

	"github.com/joho/godotenv"
)

const version = "1.0.0"

type Celeritas struct {
	AppName  string
	Debug    bool
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	RootPath string
	config   config
	Routes   *chi.Mux
	Render   *render.Render
	JetViews *jet.Set
}

type config struct {
	port     string
	renderer string
}

func New(rootPath string) (*Celeritas, error) {
	pathConfig := initPath{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "public", "logs", "tmp", "data", "middleware"},
	}

	c := &Celeritas{}
	err := c.init(pathConfig)
	if err != nil {
		return nil, err
	}
	infoLog, errLog := c.startLoggers()
	c.InfoLog = infoLog
	c.ErrorLog = errLog
	c.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	c.Version = version
	c.RootPath = rootPath
	c.Routes = c.routes().(*chi.Mux)

	c.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"), // get it from config file of the app
	}
	var views = jet.NewSet(
		jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", c.RootPath)),
		jet.InDevelopmentMode(),
	)
	c.JetViews = views

	c.createRenderer()

	return c, nil
}

func (c *Celeritas) ListenAndServe() {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", c.config.port),
		ErrorLog:     c.ErrorLog,
		Handler:      c.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	c.InfoLog.Printf("Starting server on port %s", c.config.port)

	err := server.ListenAndServe()
	c.ErrorLog.Fatal(err)
}

func (c *Celeritas) init(p initPath) error {
	root := p.rootPath

	for _, folder := range p.folderNames {
		err := c.CreateDirIfNotExist(root + "/" + folder)
		if err != nil {
			return err
		}
	}

	err := godotenv.Load(root + "/.env")
	if err != nil {
		log.Default().Println("Error loading .env file")
	}

	return nil
}

func (c *Celeritas) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (c *Celeritas) createRenderer() {
	var r render.Render
	r.Renderer = c.config.renderer
	r.RootPath = c.RootPath
	r.Port, _ = strconv.Atoi(c.config.port)
	r.Secure = false // todo: get it from config file of the app
	r.ServerName = c.AppName
	r.JetViews = c.JetViews

	c.Render = &r
}