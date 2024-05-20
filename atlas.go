package atlas

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/younesi/atlas/cache"
	"github.com/younesi/atlas/session"

	"github.com/CloudyKit/jet/v6"

	"github.com/go-chi/chi/v5"
	"github.com/younesi/atlas/render"

	"github.com/joho/godotenv"

	"github.com/alexedwards/scs/v2"
)

const version = "1.0.0"

type Atlas struct {
	AppName       string
	Debug         bool
	Version       string
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	RootPath      string
	Routes        *chi.Mux
	Render        *render.Render
	Session       *scs.SessionManager
	DB            Database
	EncryptionKey string
	Cache         cache.Cache
	config        config
}

type config struct {
	port        string
	renderer    string
	cookie      cookieConfig
	sessionType string
	database    databaseConfig
	redis       redisConfig
}

func New(rootPath string) (*Atlas, error) {
	pathConfig := initPath{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "public", "logs", "tmp", "data", "middleware"},
	}

	a := &Atlas{}
	err := a.init(pathConfig)
	if err != nil {
		return nil, err
	}
	infoLog, errLog := a.startLoggers()
	a.InfoLog = infoLog
	a.ErrorLog = errLog

	dbConfig := databaseConfig{
		dsn:      a.BuildDSN(),
		database: os.Getenv("DATABASE_TYPE"),
	}
	// connect to DB
	if dbConfig.database != "" {
		db, err := a.OpenDB(dbConfig.database, dbConfig.dsn)
		if err != nil {
			a.ErrorLog.Println(err)
			os.Exit(1)
		}

		a.DB = Database{
			Type: os.Getenv("DATABASE_TYPE"),
			Pool: db,
		}
	}

	a.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	a.Version = version
	a.RootPath = rootPath

	redisConfig := redisConfig{
		host:     os.Getenv("REDIS_HOST"),
		password: os.Getenv("REDIS_PASSWORD"),
		prefix:   os.Getenv("REDIS_PREFIX"),
	}

	a.config = config{
		port:     os.Getenv("PORT"),
		renderer: os.Getenv("RENDERER"), // get it from config file of the app
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifetime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSIST"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		sessionType: os.Getenv("SESSION_TYPE"),
		database:    dbConfig,
		redis:       redisConfig,
	}

	// create session
	sess := session.Session{
		CookieLifetime: a.config.cookie.lifetime,
		CookiePersist:  a.config.cookie.persist,
		CookieName:     a.config.cookie.name,
		CookieSecure:   a.config.cookie.secure,
		CookieDomain:   a.config.cookie.domain,
		SessionType:    a.config.sessionType,
		DBPool:         a.DB.Pool,
	}

	sss := sess.InitSession()
	a.Session = sss
	a.EncryptionKey = os.Getenv("KEY")

	if os.Getenv("CACHE") == strings.ToLower("redis") {

		myRedisCache := a.createRedisCacheClient()
		a.Cache = myRedisCache
	}

	var views = jet.NewSet(
		jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", a.RootPath)),
		jet.InDevelopmentMode(),
	)
	a.createRenderer(views)

	a.Routes = a.routes().(*chi.Mux)

	return a, nil
}

func (a *Atlas) ListenAndServe() {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", a.config.port),
		ErrorLog:     a.ErrorLog,
		Handler:      a.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	defer a.DB.Pool.Close()

	a.InfoLog.Printf("Starting server on port %s", a.config.port)

	err := server.ListenAndServe()
	a.ErrorLog.Fatal(err)
}

func (a *Atlas) init(p initPath) error {
	root := p.rootPath

	for _, folder := range p.folderNames {
		err := a.CreateDirIfNotExist(root + "/" + folder)
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

func (a *Atlas) startLoggers() (*log.Logger, *log.Logger) {
	var infoLog *log.Logger
	var errorLog *log.Logger

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (a *Atlas) createRenderer(jetViews *jet.Set) {
	var r render.Render
	r.Renderer = a.config.renderer
	r.RootPath = a.RootPath
	r.Port, _ = strconv.Atoi(a.config.port)
	r.Secure = false // todo: get it from config file of the app
	r.ServerName = a.AppName
	r.JetViews = jetViews
	r.Session = a.Session

	a.Render = &r
}

func (a *Atlas) BuildDSN() string {
	var dsn string

	switch os.Getenv("DATABASE_TYPE") {
	case "postgres", "postgresql":
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			os.Getenv("DATABASE_HOST"),
			os.Getenv("DATABASE_PORT"),
			os.Getenv("DATABASE_USER"),
			os.Getenv("DATABASE_NAME"),
			os.Getenv("DATABASE_SSL_MODE"),
		)

		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("%s password=%s", dsn, os.Getenv("DATABASE_PASS"))
		}
	default:
		a.ErrorLog.Println("could not create DB  DSN")
	}

	return dsn
}

func (a *Atlas) createRedisCacheClient() *cache.RedisCache {
	return &cache.RedisCache{
		Conn:   a.createRedisPool(),
		Prefix: a.config.redis.prefix,
	}
}

func (a *Atlas) createRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   100,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",
				a.config.redis.host,
				redis.DialPassword(a.config.redis.password))
		},
	}
}
