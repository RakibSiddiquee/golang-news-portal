package main

import (
	"database/sql"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type application struct {
	appName string
	server  server
	debug   bool
	errLog  *log.Logger
	infoLog *log.Logger
	view    *jet.Set
	session *scs.SessionManager
}

type server struct {
	host string
	port string
	url  string
}

func main() {
	server := server{
		host: "localhost",
		port: "8090",
		url:  "http://localhost:8090",
	}

	db, err := openDB("")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &application{
		server:  server,
		appName: "Golang News Portal",
		debug:   true,
		infoLog: log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate|log.Lshortfile),
		errLog:  log.New(os.Stdout, "ERROR\t", log.Ltime|log.Ldate|log.Llongfile),
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

	if err := app.listenAndServer(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hello world!")
}

// openDB is used to open database
func openDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dns)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
