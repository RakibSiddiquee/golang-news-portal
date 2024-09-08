package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/RakibSiddiquee/golang-news-portal/models"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type application struct {
	appName string
	server  server
	debug   bool
	errLog  *log.Logger
	infoLog *log.Logger
	view    *jet.Set
	session *scs.SessionManager
	Models  models.Models
}

type server struct {
	host string
	port string
	url  string
}

func main() {
	migrate := flag.Bool("migrate", false, "should migrate - drop all tables")
	flag.Parse()

	// Init server
	server := server{
		host: "localhost",
		port: "8090",
		url:  "http://localhost:8090",
	}

	// Open database connection
	db2, err := openDB("postgres://postgres:postgres@localhost/golang_news_portal?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db2.Close()

	// Init upper/db
	upper, err := postgresql.New(db2)
	if err != nil {
		log.Fatal(err)
	}
	defer func(upper db.Session) {
		err := upper.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(upper)

	// run migration
	if *migrate {
		fmt.Println("Running migrations...")
		err = runMigrate(upper)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Migrations complete")
	}

	// Init application
	app := &application{
		server:  server,
		appName: "Golang News Portal",
		debug:   true,
		infoLog: log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate|log.Lshortfile),
		errLog:  log.New(os.Stdout, "ERROR\t", log.Ltime|log.Ldate|log.Llongfile),
		Models:  models.New(upper),
	}

	// Init jet template
	if app.debug {
		app.view = jet.NewSet(jet.NewOSFileSystemLoader("./views"), jet.InDevelopmentMode())
	} else {
		app.view = jet.NewSet(jet.NewOSFileSystemLoader("./views"))
	}

	// Init session
	app.session = scs.New()
	app.session.Lifetime = 24 * time.Hour
	app.session.Cookie.Persist = true
	app.session.Cookie.Name = strings.ReplaceAll(app.appName, " ", "-")
	app.session.Cookie.Domain = app.server.host
	app.session.Cookie.SameSite = http.SameSiteStrictMode
	app.session.Store = postgresstore.New(db2)

	if err := app.listenAndServer(); err != nil {
		log.Fatal(err)
	}
}

// openDB is used to open database
func openDB(dsn string) (*sql.DB, error) {
	odb, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = odb.Ping(); err != nil {
		return nil, err
	}
	return odb, nil
}

// runMigrate is used to run migration files
func runMigrate(db db.Session) error {
	script, err := os.ReadFile("./migrations/tables.sql")
	if err != nil {
		return err
	}

	_, err = db.SQL().Exec(string(script))

	return err
}
